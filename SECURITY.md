# Security Policy

## Supported versions

Security fixes are applied to the `main` branch only. There are no long-lived
release branches or published package channels for this repository.

## Reporting a vulnerability

Please report security issues privately through GitHub Security Advisories
(Security → Advisories → New draft security advisory) on this repository.

Do not open a public issue for vulnerabilities that could be exploited by
untrusted User-Agent or Client Hints input.

Include:

- A clear description of the issue and impact
- Steps to reproduce (minimal UA string, headers, or Go snippet when possible)
- Affected commit or module version if known

You can expect an acknowledgment when the report is reviewed. Fixes will be
coordinated privately when appropriate, then disclosed after a patch is on
`main`.
