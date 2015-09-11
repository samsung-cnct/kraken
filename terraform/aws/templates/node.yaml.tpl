#cloud-config

---
write_files:
  - path: /etc/inventory.ansible
    content: |
      [nodes]
      ${short_name} ansible_ssh_host=$private_ipv4

      [nodes:vars]
      ansible_connection=ssh
      ansible_python_interpreter="PATH=/home/core/bin:$PATH python"
      ansible_ssh_user=core
      ansible_ssh_private_key_file=/opt/ansible/private_key
      cluster_master_record=${cluster_master_record}
      cluster_proxy_record=${cluster_proxy_record}
      cluster_name=${cluster_name}
      dns_domain=${dns_domain}
      dns_ip=${dns_ip}
      dockercfg_base64=${dockercfg_base64}
      etcd_private_ip=${etcd_private_ip}
      etcd_public_ip=${etcd_public_ip}
      hyperkube_deployment_mode=${hyperkube_deployment_mode}
      hyperkube_image=${hyperkube_image}
      interface_name=${interface_name}
      kraken_services_branch=${kraken_services_branch}
      kraken_services_dirs=${kraken_services_dirs}
      kraken_services_repo=${kraken_services_repo}
      kube_apiserver_v=${kube_apiserver_v}
      kube_controller_manager_v=${kube_controller_manager_v}
      kube_proxy_v=${kube_proxy_v}
      kube_scheduler_v=${kube_scheduler_v}
      kubelet_v=${kubelet_v}
      kubernetes_api_version=${kubernetes_api_version}
      kubernetes_binaries_uri=${kubernetes_binaries_uri}
      logentries_token=${logentries_token}
      logentries_url=${logentries_url}
      master_private_ip=${master_private_ip}
      master_public_ip=${master_public_ip}
coreos:
  etcd2:
    proxy: on
    listen-client-urls: http://0.0.0.0:2379,http://0.0.0.0:4001
    advertise-client-urls: http://0.0.0.0:2379,http://0.0.0.0:4001
    initial-cluster: etcd=http://${etcd_private_ip}:2380
  fleet:
    etcd-servers: http://$private_ipv4:4001
    public-ip: $private_ipv4
    metadata: "role=node"
  flannel:
    etcd-endpoints: http://${etcd_private_ip}:4001
    interface: $private_ipv4
  units:
    - name: format-storage.service
      command: start
      content: |
        [Unit]
        Description=Formats a drive
        [Service]
        Type=oneshot
        RemainAfterExit=yes
        ExecStart=/usr/sbin/wipefs -f ${format_docker_storage_mnt}
        ExecStart=/usr/sbin/mkfs.ext4 -F ${format_docker_storage_mnt}
    - name: var-lib-docker.mount
      command: start
      content: |
        [Unit]
        Description=Mount to /var/lib/docker
        Requires=format-storage.service
        After=format-storage.service
        Before=docker.service
        [Mount]
        What=${format_docker_storage_mnt}
        Where=/var/lib/docker
        Type=ext4
    - name: docker-tcp.socket
      command: start
      enable: true
      content: |
        [Unit]
        Description=Docker Socket for the API

        [Socket]
        ListenStream=0.0.0.0:4243
        BindIPv6Only=both
        Service=docker.service

        [Install]
        WantedBy=sockets.target
    - name: setup-network-environment.service
      command: start
      content: |
        [Unit]
        Description=Setup Network Environment
        Documentation=https://github.com/kelseyhightower/setup-network-environment
        Requires=network-online.target
        After=network-online.target
        Before=flanneld.service

        [Service]
        ExecStartPre=-/usr/bin/mkdir -p /opt/bin
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://github.com/kelseyhightower/setup-network-environment/releases/download/v1.0.0/setup-network-environment
        ExecStartPre=/usr/bin/chmod +x /opt/bin/setup-network-environment
        ExecStart=/opt/bin/setup-network-environment
        RemainAfterExit=yes
        Type=oneshot
    - name: fleet.service
      command: start
    - name: etcd2.service
      command: start
    - name: flanneld.service
      command: start
    - name: systemd-journal-gatewayd.socket
      command: start
      enable: yes
      content: |
        [Unit]
        Description=Journal Gateway Service Socket
        [Socket]
        ListenStream=/var/run/journald.sock
        Service=systemd-journal-gatewayd.service
        [Install]
        WantedBy=sockets.target
    - name: generate-ansible-keys.service
      command: start
      content: |
        [Unit]
        Description=Generates SSH keys for ansible container
        [Service]
        Type=oneshot
        RemainAfterExit=yes
        ExecStart=/usr/bin/bash -c "ssh-keygen -f /home/core/.ssh/ansible_rsa -N ''"
        ExecStart=/usr/bin/bash -c "cat /home/core/.ssh/ansible_rsa.pub >> /home/core/.ssh/authorized_keys"
    - name: kraken-git-pull.service
      command: start
      content: |
        [Unit]
        Requires=generate-ansible-keys.service
        After=generate-ansible-keys.service
        Description=Fetches kraken repo
        [Service]
        Type=oneshot
        RemainAfterExit=yes
        ExecStart=/usr/bin/rm -rf /opt/kraken
        ExecStart=/usr/bin/git clone -b ${kraken_branch} ${kraken_repo} /opt/kraken
    - name: write-sha-file.service
      command: start
      content: |
        [Unit]
        Requires=kraken-git-pull.service
        After=kraken-git-pull.service
        Description=writes optional sha to a file
        [Service]
        Type=oneshot
        RemainAfterExit=yes
        ExecStart=/usr/bin/echo ${kraken_commit} > /opt/kraken/commit.sha
    - name: fetch-kraken-commit.service
      command: start
      content: |
        [Unit]
        Requires=write-sha-file.service
        After=write-sha-file.service
        Description=fetches an optional commit
        ConditionFileNotEmpty=/opt/kraken/commit.sha
        [Service]
        Type=oneshot
        RemainAfterExit=yes
        WorkingDirectory=/opt/kraken
        ExecStart=/usr/bin/git fetch ${kraken_repo} +refs/pull/*:refs/remotes/origin/pr/*
        ExecStart=/usr/bin/git checkout -f ${kraken_commit}
    - name: ansible-in-docker.service
      command: start
      content: |
        [Unit]
        Requires=write-sha-file.service
        After=fetch-kraken-commit.service
        Description=Runs a prebaked ansible container
        [Service]
        Type=simple
        RemainAfterExit=yes
        Restart=on-failure
        RestartSec=3
        ExecStart=/usr/bin/docker run -v /etc/inventory.ansible:/etc/inventory.ansible -v /opt/kraken:/opt/kraken -v /home/core/.ssh/ansible_rsa:/opt/ansible/private_key -v /var/run:/ansible -e ANSIBLE_HOST_KEY_CHECKING=False -e quay.io/samsung_ag/kraken_ansible /sbin/my_init --skip-startup-files --skip-runit -- ansible-playbook /opt/kraken/ansible/iaas_provision.yaml -i /etc/inventory.ansible
  update:
    group: ${coreos_update_channel}
    reboot-strategy: ${coreos_reboot_strategy}
