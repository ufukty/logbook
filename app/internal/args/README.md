# args

Values of config parameters might depend on multiple parameters. Such as the service binary represents or deployment environment such as staging, production or local. Also some parameters might need to have same value for different environments or services.

Args expects all binaries need at least 3 config files for:

-   service-dependent parameters `-s`,
-   environment-dependent parameters `-d`,
-   independent parameters `-i`
