# Problems with legacy codebase that needs to be rectified
**Below**, we outline the primary issues identified in the legacy codebase that need to be addressed. This list highlights key areas of concern but is not exhaustive; the team will work to resolve these and any additional problems that arise.

## SQL Injection 
SQL Injection is a serious issue that can be exploited by a malicious actor. The team will strive to ensure the use of placeholder variables in SQL strings within the rewrite that disallows SQL injection. 

## Password Encryption
The legacy codebase hashes user passwords with MD5 which is [**not** a secure hashing algorithm.](https://medium.com/@shital.pimpale5/why-is-md5-no-longer-secure-4c2c89713eb2)

## Hardcoded Environment Variables
The legacy project contains hardcoded environment variables which is not beneficial for the team. Hardcoded environment variables like API keys may be exploited by bad actors which will reduce the efficiency of the team. 