FROM centos:7

WORKDIR /root

RUN                                                           \
    yum install -y epel-release &&                            \
    yum update -y &&                                          \
    yum install -y wget &&                                    \
    for f in                                                  \
        ceph-14.2.4-0.el7.x86_64.rpm                          \
        ceph-base-14.2.4-0.el7.x86_64.rpm                     \
        ceph-common-14.2.4-0.el7.x86_64.rpm                   \
        ceph-fuse-14.2.4-0.el7.x86_64.rpm                     \
        ceph-mds-14.2.4-0.el7.x86_64.rpm                      \
        ceph-mgr-14.2.4-0.el7.x86_64.rpm                      \
        ceph-mon-14.2.4-0.el7.x86_64.rpm                      \
        ceph-osd-14.2.4-0.el7.x86_64.rpm                      \
        ceph-radosgw-14.2.4-0.el7.x86_64.rpm                  \
        ceph-resource-agents-14.2.4-0.el7.x86_64.rpm          \
        ceph-selinux-14.2.4-0.el7.x86_64.rpm                  \
        ceph-test-14.2.4-0.el7.x86_64.rpm                     \
        cephfs-java-14.2.4-0.el7.x86_64.rpm                   \
        libcephfs-devel-14.2.4-0.el7.x86_64.rpm               \
        libcephfs2-14.2.4-0.el7.x86_64.rpm                    \
        libcephfs_jni-devel-14.2.4-0.el7.x86_64.rpm           \
        libcephfs_jni1-14.2.4-0.el7.x86_64.rpm                \
        librados-devel-14.2.4-0.el7.x86_64.rpm                \
        librados2-14.2.4-0.el7.x86_64.rpm                     \
        libradospp-devel-14.2.4-0.el7.x86_64.rpm              \
        libradosstriper-devel-14.2.4-0.el7.x86_64.rpm         \
        libradosstriper1-14.2.4-0.el7.x86_64.rpm              \
        librbd-devel-14.2.4-0.el7.x86_64.rpm                  \
        librbd1-14.2.4-0.el7.x86_64.rpm                       \
        librgw-devel-14.2.4-0.el7.x86_64.rpm                  \
        librgw2-14.2.4-0.el7.x86_64.rpm                       \
        python-ceph-argparse-14.2.4-0.el7.x86_64.rpm          \
        python-ceph-compat-14.2.4-0.el7.x86_64.rpm            \
        python-cephfs-14.2.4-0.el7.x86_64.rpm                 \
        python-rados-14.2.4-0.el7.x86_64.rpm                  \
        python-rbd-14.2.4-0.el7.x86_64.rpm                    \
        python-rgw-14.2.4-0.el7.x86_64.rpm                    \
        python36-ceph-argparse-14.2.4-0.el7.x86_64.rpm        \
        python36-cephfs-14.2.4-0.el7.x86_64.rpm               \
        python36-rados-14.2.4-0.el7.x86_64.rpm                \
        python36-rbd-14.2.4-0.el7.x86_64.rpm                  \
        python36-rgw-14.2.4-0.el7.x86_64.rpm                  \
        rados-objclass-devel-14.2.4-0.el7.x86_64.rpm          \
        rbd-fuse-14.2.4-0.el7.x86_64.rpm                      \
        rbd-mirror-14.2.4-0.el7.x86_64.rpm                    \
        rbd-nbd-14.2.4-0.el7.x86_64.rpm                       \
    ;do                                                       \
        wget http://hk.ceph.com/rpm-nautilus/el7/x86_64/$f    \
    ;done           &&                                        \
    yum install -y *.rpm  &&                                  \
    rm -f *.rpm     &&                                        \
    yum clean all   &&                                        \
    rpmdb --rebuilddb
