# Yggdrasil Platform

## Mission Statement

The Yggdrasil Platform is intended to be a self-managing
high-availability infrastructure platform that can be easily
bootstrapped on any collection of servers, providing the
tools necessary for developers targeting the platform to
design and build applications that are secure by default.

This security will descend from the security of the platform
itself, which will be designed with a zero-trust infastructure.

## Architecture Design

The platform shall consist of two distinct top-level components:
storage management and workload management. Workload management
will be accomplished by Rancher RKE2 and Storage management will
be accomplished using Ceph. These two will be linked using the
Persistent Volume facilities built into Kubernetes and a dynamic
provisioner for allocating PVCs as needed.

All nodes for the Workload management system will boot over the
network using iPXE to load signed and verify cryptographically-signed images.
A continuous process will update these images periodically by installing
updates from the OS vendor as well as any updates published by Rancher
to RKE2. Users are expected to track upstream deprecation cycles, as features
will be removed by this process when they are removed in the upstream
services.

Updates will be applied to the cluster via periodic reboots of each node.
This will be done without warning to any workload, which will enforce
a failure-tolerant design strategy for any application targeting the platform.

On boot, each node will be provided with the most recent image which is not
known to be problematic, which shall be a control plane or an agent node
depending on the identified needs of the workload management system. In
particular, a control plane image will be provided if there are insufficiently
many control plane nodes in the cluster; otherwise an agent image will be
supplied to the machine.

This image server will be hosted as a workload running on the workload
management system itself.
