# Nidavellir

Nidavellir is the Yggrasil subsystem responsible for building the boot images
for Yggdrasil nodes.

When possible, the images themselves shall be unmodified mirrors of the upstream
Alpine initfs and kernel files, with any runtime changes made in the form
of Alpine apkovl files that can be applied over the base OS at boot.


