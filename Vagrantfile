Vagrant.configure("2") do |config|
  nodes = {
    "k8s-master" => "192.168.56.10",
    "k8s-worker" => "192.168.56.11"
  }

  config.vm.box = "ubuntu/jammy64"

  nodes.each do |name, ip|
    config.vm.define name do |node|
      node.vm.hostname = name
      node.vm.network "private_network", ip: ip
      node.vm.provider "virtualbox" do |vb|
        vb.memory = 4096
        vb.cpus = 2
      end
      node.vm.provision "shell", inline: <<-SHELL
        sudo apt-get update
        sudo apt-get install -y curl vim net-tools
      SHELL
    end
  end
end
