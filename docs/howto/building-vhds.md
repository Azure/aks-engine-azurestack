# How to build VHDs

## Windows VHD

To build the Windows VHD

1. Enter our dev container

    ```bash
    make dev
    ```

1. Set the following environment varibles

    |   |   |
    |---|---|
    | CLIENT_ID | Azure service princicle ID |
    | CLIENT_SECRET | Azure service principle secret |
    | TENANT_ID | Azure tennant ID |
    |  AZURE_RESOURCE_GROUP_NAME | Resource group name to use / create |
    | AZURE_LOCATION | Azure region to use |
    | AZURE_VM_SIZE | VM size packer will use (Standard_D4s_v3 is recommended) |
    | CONTAINER_RUNTIME | Container runtime VHD will be generated for (Docker or Containerd) |
    | WINDOWS_SERVER_VERSION | Windows Server OS version vHD will be generated for (2019 or 2004) |

1. Run packer

   ```bash
   make run-packer-windows
   ```

Once complete, packer will output a URI to the VHD it created.

## Linux VHD

TODO

## Temporary eMudhra remediation in Azure Image Builder

The [Linux AIB template](../../vhd/aib/AKSeImageTemplateLinux.json) and
[Windows AIB template](../../vhd/aib/AKSeImageTemplateWindows.json) include
build-time remediation and a final verification gate. These steps do not change
the source image references or remediate existing deployed nodes. The legacy
Packer templates above are not wired to these AIB steps.

Set `github_branch_name` to a published branch or commit containing the updated
templates and remediation scripts. A local change alone is not available to the
image builder's `raw.githubusercontent.com` downloads.

### Ubuntu 22.04

The [Linux script](../../vhd/packer/remediate-emudhra.sh) updates APT metadata,
installs or upgrades `ca-certificates` and every installed architecture of
`libnss3`, and regenerates the system CA bundle. It requires fully installed
packages at or above these Jammy-specific versions:

| Package | Minimum fixed version |
|---|---|
| `ca-certificates` | `20260816ubuntu0.22.04.1` |
| `libnss3` | `2:3.98-0ubuntu0.22.04.5` |

Version comparisons use `dpkg`, including the NSS epoch. Verification also parses
the system CA bundle and rejects the withdrawn emSign/AAA root names if they
were reintroduced locally. This name check is supplementary to the vendor
package checks, not a general-purpose SPKI blocklist.

The remediation runs before dependency installation. Verification runs again
after hardening and cleanup. An APT failure, missing or old package, incomplete
package configuration, unreadable bundle, or remaining withdrawn root fails the
build. Internet access to the configured Ubuntu package repositories is required.
The remediation shell script is stored and checked out with LF line endings,
enforced by `.gitattributes`, so Windows-authored changes remain executable on
the Linux builder.

Sources: [Canonical's advisory](https://ubuntu.com/blog/emudhra-ca-compromise-fixes-available),
[CA update](https://ubuntu.com/security/notices/USN-8944-1),
and [NSS update](https://ubuntu.com/security/notices/USN-8947-1).

### Windows Server 2019

The [Windows script](../../vhd/packer/remediate-emudhra.ps1) uses Microsoft's
documented forced CTL refresh:

```text
certutil -f -syncWithWU <ctl-data-directory>
certutil -f -verifyCTL Disallowed <ctl-data-directory> <untrusted-probe.cer>
```

Synchronization supplies `disallowedcert.sst`, which `verifyCTL` uses to locate
certificates matching CTL entries. The metadata is retained under
`C:\akse-cache\emudhra-ctl` for the final and offline verification steps;
`-CtlDirectory` can override this location. Synchronization alone is not proof
that Windows installed the CTL.

The self-signed probe is generated in memory, exported without its private key,
and never installed into a certificate store. This supported operation refreshes
**both AuthRoot and Disallowed CTLs**. It does not manually block parent root CAs,
but normal Microsoft trust-program changes can affect certificate compatibility.
Windows Update CTL endpoints must be reachable and CTL updates must not be
blocked by policy.

The synthetic probe deliberately does not belong to the Disallowed CTL, so the
forced-refresh command can return
`0x80094005` (`CERTSRV_E_INVALID_CA_CERTIFICATE`) for that unmatched probe.
Only that specific status is accepted for the probe command. It is logged as
an expected non-match, not treated as proof that the refresh succeeded.

The build then requires successful
`certutil -verifyCTL Disallowed <ctl-data-directory>` verification, without the probe,
of the local registry-cached CTL. The sequence number is decoded from the
persisted `DisallowedCertEncodedCtl` using Windows CryptoAPI, not by parsing
localized command output. It must equal:

```text
01dd5f6ca45c9d86
```

[Microsoft's incident announcement](https://github.com/TrustedRootProgram/Program-Requirements/blob/main/Announcements.md#october-26-2026)
identifies that CTL as blocking all eight compromised keys by SPKI, including
the identified cross-signed certificate paths. A different sequence, even a
numerically newer one, fails closed until its mitigation coverage is reviewed
and the expected value is explicitly updated.

The final customizer repeats local verification after the restarts and explicit
sysprep preparation steps, writing
`C:\buildArtifacts\emudhra-remediation.txt`. It does not treat a downloaded CTL,
an SST import, or an empty intermediate certificate cache as proof of
remediation. Synchronization errors, other probe errors, and all final
verification errors are fatal.

See Microsoft's [certutil reference](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/certutil#-verifyctl)
and [CTL servicing guidance](https://learn.microsoft.com/en-us/windows-server/identity/ad-cs/configure-trusted-roots-disallowed-certificates).

### Qualification before publishing

Run an actual AIB build on each target OS. After capture/generalization, boot a
VM from each resulting image with network access disabled and run the relevant
script's verification-only mode on that VM:

```bash
sudo bash /home/packer/remediate-emudhra.sh --verify-only
```

```powershell
& C:\akse-cache\remediate-emudhra.ps1 -VerifyOnly
```

Confirm that the evidence survives capture and that affected certificate paths
produce explicit distrust in Windows chain-validation tests. An expired
certificate, missing issuer, or generic TLS failure is not proof of explicit
distrust. The build-time CTL identity check does not replace these image
qualification tests.

These steps cover the OS trust stores only. Application-private CA bundles,
Java stores, cached container images, and statically linked executables require
separate remediation. In particular, rebuilding the AKS Engine binary with the
patched Go toolchain remains necessary.