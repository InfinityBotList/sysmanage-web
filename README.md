# github.com/infinitybotlist/sysmanage-web

**See the ``example`` directory as a example on how to use sysmanage-web library to configure a proper executable**

# Getting started

1. Copy the contents of the ``example`` directory out to a seperate folder (``sysmanage`` etc.)
2. Run ``git init`` and ``go mod tidy`` to download any deps needed
3. Run ``npm i`` in ``frontend``.
4. Personalize it to your preferences, and build using ``make``

Allows management of our systems, though it can be used by anyone. 
Core plugins included by default are ``nginx``, ``systemd``, ``persist``, ``frontend`` and ``actions``.

This should be running under ``deployproxy`` or some other authentication proxy/system for additional security. If you wish to setup a different authentication proxy or do not want deployproxy auth checks, such as when performing initial bootstrapping, you can set ``dp_disable`` in ``config.yaml``.
