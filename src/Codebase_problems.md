# Problems with the codebase
**Problems:**

* SQL Injection ,linje 147 

* Bruger passwords hashes med MD5, https://medium.com/@shital.pimpale5/why-is-md5-no-longer-secure-4c2c89713eb2  
  * Usaltede passwords?  

* SECRET_KEY = 'development key' hardcode er ikke godt sikkerhedsmæssigt, særligt hvis det eksponeres til Github ved push. 