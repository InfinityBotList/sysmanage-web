# sysmanage-web

*All rights reserved...*

Allows management of our systemd services, as well as creating/deleting them, quickly through the browser. This replaces ``infra-scripts`` and ``service-gen``.

**This must be running under deployproxy for additional security. When performing initial bootstrapping however, you can set dp_disable in ``config.yaml`` to disable this**

1. Build ``frontend`` first using its README.md
2. Run ``make`` to build the backend after creating the ``config.yaml`` file.  
 
# Rust

```yaml
git:
    repo: https://github.com/infinitybotlist/persepolis
    ref: refs/heads/main
    build_commands:
        - /root/.cargo/bin/cargo build --release
        - systemctl stop persepolis
        - rm -vf persepolis
        - mv -vf target/release/persepolis .
        - systemctl start persepolis
    env:
        DATABASE_URL: postgres:///infinity
        RUSTFLAGS: -C target-cpu=native -C link-arg=-fuse-ld=lld
    allow_dirty: true
    config_files:
        - config.yaml
```

Use the above template for the git integration for rust