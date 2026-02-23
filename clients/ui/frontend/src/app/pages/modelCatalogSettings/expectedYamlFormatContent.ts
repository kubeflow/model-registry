export const EXPECTED_YAML_FORMAT_CONTENT = `all:
  children:
    control_nodes:
      hosts:
        node1:
          custom_hostname: <VALUE>
          management_network:
            network1:
              ip: <VALUE>
        node2:
          custom_hostname: <VALUE>
          management_network:
            network1:
              ip: <VALUE>
        node3:
          custom_hostname: <VALUE>
          management_network:
            network1:
              ip: <VALUE>
    switches:
    #BEGIN BGP
      vars:
        cp4d_asplain: <VALUE>
        cp4d_network: <VALUE>
        cp4d_network_vip: <VALUE>
      hosts:
        FabSw1a:
          ansible_host: localhost
          vrr_ip_addr: <VALUE>
          cp4d_routerID: 9.0.62.1
          isl_peer: 9.0.255.2
          bgp_links:
              link1:
              swp: <VALUE>
              neighbor: <VALUE>
              ip_addr: <VALUE>
              mtu: 9000
              link_speed: 10000
        FabSw1b:
          ansible_host: localhost
          vrr_ip_addr: <VALUE>
          cp4d_routerID: 9.0.62.2
          isl_peer: 9.0.255.1
          bgp_links:
            link1:
              swp: <VALUE>
              neighbor: <VALUE>
              ip_addr: <VALUE>
              mtu: 9000
              link_speed: 10000
    #END BGP

    #BEGIN L2
    switches:
      hosts:
        FabSw1a:
          ansible_host: localhost
          external_connection_config:
            external_link1:
              switch_ports: ['<VALUE>', '<VALUE>']
              port_config:
                mtu: 9000
                link_speed: 10000
              vlans: ['VALUE']
              strict_vlan: <VALUE>
              name: <VALUE>
              lacp_link: True
              lacp_rate: Fast
              clag_id: 100
              partner_switch: 'FabSw1b'
    #END L2

  vars:
    app_fqdn: <VALUE>
    #(pick from timedatectl list-timezones), default is EDT
    timezone: "<OPTIONAL>"
    #must begin with server or pool
    time_servers: ["<OPTIONAL>"]
    dns_servers: ["<VALUE>"]
    dns_search_strings: ["<OPTIONAL>"]
    smtp_servers: ["<OPTIONAL>"]
    management_network:
      network1:
        subnet: <VALUE>
        # just number, no slash 
        prefix: <VALUE>
        gateway: <VALUE>
        floating_ip: <VALUE>
        mtu: <OPTIONAL>
        custom_routes: <OPTIONAL>
    application_network_enabled: False
    openshift_networking_enabled: False
    policy_based_routing_enabled: True
    application_network:
      network1:
        default_gateway: true
        vlan: <VALUE>
        # just number, no slash 
        prefix: <VALUE>
        gateway: <VALUE>
        floating_ip: <VALUE>
        mtu: <OPTIONAL>
        custom_routes: <OPTIONAL>
        additional_openshift_ipaddrs: ["<OPTIONAL>"]
        additional_openshift_routes: ["<OPTIONAL>"]
`;
