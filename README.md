# sysmanage-web

*All rights reserved*

Allows management of our systemd services, as well as creating/deleting them, quickly through the browser. This replaces ``infra-scripts`` and ``service-gen``.

**This must be running under deployproxy for additional security. When performing initial bootstrapping however, you can set dp_disable in ``config.yaml`` to disable this**

1. Build ``frontend`` first using its README.md
2. Run ``make`` to build the backend after creating the ``config.yaml`` file