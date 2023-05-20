# github.com/infinitybotlist/sysmanage-web

**See ``sysmanage-ibl`` as a example on how to use sysmanage-web**

Allows management of our systems, though it can be used by anyone. 
Core plugins included by default are ``nginx``, ``systemd``, ``persist``, ``frontend`` and ``actions``.

This should be running under ``deployproxy`` or some other authentication proxy/system for additional security. If you wish to setup a different authentication proxy or do not want deployproxy auth checks, such as when performing initial bootstrapping, you can set ``dp_disable`` in ``config.yaml``.
