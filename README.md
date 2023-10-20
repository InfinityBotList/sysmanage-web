# Sysmanage-web

Simple system management framework for small-medium hosts/setups

**See the ``example`` directory as a example on how to use sysmanage-web library to configure a proper executable**

# Getting started

## Automatic

1. Run ``projbuilder/build.py``

## Manual

1. Copy the contents of the ``example`` directory out to a seperate folder (``sysmanage`` etc.)
2. Run ``git init`` and ``go mod tidy`` to download any deps needed
3. Run ``npm i`` in ``frontend``.
4. Personalize it to your preferences, and build using ``make``

Allows management of our systems, though it can be used by anyone. 
Core plugins included by default are ``nginx``, ``systemd``, ``persist``, ``frontend`` and ``actions``.

This should be running under ``deployproxy`` or some other authentication proxy/system for additional security. If you wish to setup a different authentication proxy or do not want deployproxy auth checks, such as when performing initial bootstrapping, you can set ``dp_disable`` in ``config.yaml``.

# Plugins

## Systemd

Systemd integration for ``sysmanage`` to allow easy system management

1. Create ``data/servicegen/server.tmpl``. Add your systemd service template here. *Example:*

```toml
[Service]
Type=simple
ExecStart={{.Command}}
User={{.User}}
Group={{.Group}}
WorkingDirectory={{.Directory}}
ExecReload=/bin/kill -s HUP $MAINPID
KillMode=mixed
TimeoutStopSec=5
PrivateTmp=true
RestartSec=1
Restart=always

[Install]
WantedBy=multi-user.target

[Unit]
PartOf={{.Target}}.target
Description="{{.Description}}"
After={{.After}}.service
```

2. Create ``data/servicegen/target.tmpl``. Add your systemd target here. *Example:*

```toml
[Unit]
Description={{.Description}}

# This collection of apps should be started at boot time.
[Install]
WantedBy=multi-user.target
```

3. Create ``data/services/_meta.yaml`` with content of ``target:`` (include the ending colon).

## Wafflepaw

Wafflepaw is a simple monitoring and alert system. It can also integrate with Instatus (Metrics, Incidents) as well as with Discord. Wafflepaw should be ran from a secondary server if possible

**Wafflepaw is still a work in progress and is not ready for use**