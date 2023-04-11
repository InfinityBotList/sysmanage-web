<script lang="ts">
	import Doc from "./Doc.svelte";
</script>

<svelte:head>
	<title>About</title>
	<meta name="description" content="About this app" />
</svelte:head>

<div class="text-column">
	<h1 class="font-semibold">Documentation</h1>

	<Doc docTitle="install-nirn-proxy">
		<p>
			Due to VPS ratelimits, IBL uses the <code>nirn-proxy</code>
			available at <a class="text-blue-400 hover:text-blue-600" href="https://github.com/germanoeich/nirn-proxy">https://github.com/germanoeich/nirn-proxy</a>
		</p>
		<ul class="list-disc list-inside">
			<li>Install <code>iblcli</code> from <a class="text-blue-400 hover:text-blue-600" href="https://github.com/infinitybotlist/iblcli">https://github.com/infinitybotlist/iblcli</a></li>
			<li>Install <code>nirn-proxy</code> and copy it to <code>/usr/bin</code></li>
			<li>The <code>baseproxy</code> service should now be working properly</li>
		</ul>
	</Doc>
	<Doc docTitle="zfs-setup">
		<p>The optimum setup for ZFS is as follows. Note that the below may be out-of-date right now and is being revised</p>
		<ul class="list-disc list-inside">
			<li>ZFS pool named silverpelt created on a 100GB DO volume or similar</li>
			<li>On silverpelt, a XFS filesystem attached to a zvol named nightheart (silverpelt/nightheart) of 90-95GB</li>
			<li>Ubuntu 22.04 LTS with zfsutils-linux used (not FUSE)</li>
			<li>
				Sunbeam as ZPOOL. ARC of 5368709120 following https://bun.uptrace.dev/postgres/tuning-zfs-aws-ebs.html#zfs-config.
			</li>
			<li>
				Set data_directory in postgres to /silverpelt/sunbeam. Use ZFS native encryption here for data protection.
			</li>
		</ul>
	</Doc>
</div>
