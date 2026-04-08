-- Complementary Sample Data for UI filtering, sorting and pagination scenarios

-- Insert Complementary Assets
INSERT INTO asset (id, name, description, createdat, lastscan) VALUES
('AST-013', 'Dell PowerEdge R650 Server', 'Backup compute node used for quarterly disaster recovery drills', '2024-04-12', '2024-10-09'),
('AST-014', 'Palo Alto PA-3220 Firewall', 'Regional office perimeter firewall in high-availability pair', '2024-04-20', '2024-10-09'),
('AST-015', 'Lenovo ThinkSystem SR630', 'Analytics processing server in internal segmentation zone', '2024-05-02', '2024-10-09'),
('AST-016', 'Cisco UCS C220 M6', 'Newly provisioned bare-metal host pending first security scan', '2024-05-12', NULL),
('AST-017', 'HPE Apollo 4200 Gen10 Plus', 'High-density backup repository server for immutable snapshots', '2024-05-18', '2024-11-02'),
('AST-018', 'Juniper SRX340 Firewall', 'Branch office secure edge gateway handling SD-WAN tunnels', '2024-05-25', '2024-11-04'),
('AST-019', 'Supermicro SYS-5019C-MR', 'Lab server used for firmware regression validation', '2024-06-01', '2024-11-06'),
('AST-020', 'Dell OptiPlex 7090', 'Finance desktop used for ERP and treasury operations', '2024-06-07', '2024-11-08'),
('AST-021', 'Arista 7280R3 Switch', 'Datacenter spine switch for east-west traffic aggregation', '2024-06-14', '2024-11-10'),
('AST-022', 'Lenovo ThinkPad T14 Gen 4', 'Security operations laptop assigned to incident response', '2024-06-21', '2024-11-12'),
('AST-023', 'Fortinet FortiGate 200F', 'Regional segmentation firewall between production and staging', '2024-06-28', '2024-11-14'),
('AST-024', 'HP ZBook Fury 16 G9', 'Mobile engineering workstation waiting for baseline enrollment', '2024-07-05', NULL),
('AST-025', 'Cisco Nexus 93180YC-FX', 'Top-of-rack switch serving Kubernetes worker pools', '2024-07-12', '2024-11-18'),
('AST-026', 'Apple Mac mini M2', 'Build agent host for iOS CI pipelines in isolated VLAN', '2024-07-19', '2024-11-20'),
('AST-027', 'HPE ProLiant DL360 Gen11', 'Authentication services host for internal identity workloads', '2024-07-26', '2024-11-22'),
('AST-028', 'Dell Latitude 7440', 'Field support laptop used for customer on-site diagnostics', '2024-08-02', '2024-11-24'),
('AST-029', 'Juniper QFX5120-48Y', 'Leaf switch connecting storage and virtualization clusters', '2024-08-09', '2024-11-26'),
('AST-030', 'Supermicro AS-1115CS-TNR', 'AI inferencing node with hardened boot configuration', '2024-08-16', '2024-11-28'),
('AST-031', 'Palo Alto PA-1410 Firewall', 'Campus internet edge firewall enforcing outbound policies', '2024-08-23', '2024-11-30'),
('AST-032', 'Cisco ISR 4431', 'Remote branch WAN router not yet integrated into scan schedule', '2024-08-30', NULL),
('AST-033', 'Lenovo ThinkCentre M90q', 'Kiosk management endpoint in manufacturing floor subnet', '2024-09-06', '2024-12-04'),
('AST-034', 'Arista 7050SX3-48YC8', 'Storage fabric switch for backup replication network', '2024-09-13', '2024-12-06'),
('AST-035', 'Dell PowerStore 1200T', 'Tier-1 storage appliance for transactional workloads', '2024-09-20', '2024-12-08'),
('AST-036', 'HPE Aruba 6300M', 'Core campus switch aggregating access layer traffic', '2024-09-27', '2024-12-10'),
('AST-037', 'Fortinet FortiSwitch 424E', 'Access switch in PCI-compliant payment processing segment', '2024-10-04', '2024-12-12'),
('AST-038', 'Apple MacBook Air M3', 'Marketing laptop staged for onboarding and pending first scan', '2024-10-11', NULL),
('AST-039', 'Cisco Meraki MX95', 'Retail edge appliance for secure SD-WAN connectivity', '2024-10-18', '2024-12-16'),
('AST-040', 'Lenovo ThinkSystem ST250 V3', 'Edge compute server delivered but not yet commissioned', '2024-10-25', NULL)
ON CONFLICT (id) DO NOTHING;

-- Insert Complementary Components
INSERT INTO component (id, name, version, vendor, type, createdat, lastscan, assetid) VALUES
('CMP-039', 'Dell iDRAC9 Firmware', '5.10.10.00', 'Dell Inc.', 'BMC', '2024-04-12', '2024-10-09', 'AST-013'),
('CMP-040', 'PAN-OS System Software', '10.2.9-h8', 'Palo Alto Networks', 'Operating System', '2024-04-20', '2024-10-09', 'AST-014'),
('CMP-041', 'Lenovo XClarity Controller', '3.20', 'Lenovo', 'Management Controller', '2024-05-02', '2024-10-09', 'AST-015'),
('CMP-042', 'Cisco CIMC Firmware', '4.2(3i)', 'Cisco Systems', 'BMC', '2024-05-12', NULL, 'AST-016'),
('CMP-043', 'HPE iLO 5 Firmware', '2.95', 'HPE', 'BMC', '2024-05-18', '2024-11-02', 'AST-017'),
('CMP-044', 'Junos OS', '21.4R3-S2', 'Juniper Networks', 'Operating System', '2024-05-25', '2024-11-04', 'AST-018'),
('CMP-045', 'Supermicro IPMI Firmware', '1.86', 'Supermicro', 'BMC', '2024-06-01', '2024-11-06', 'AST-019'),
('CMP-046', 'Dell UEFI BIOS', '1.19.0', 'Dell Inc.', 'UEFI', '2024-06-07', '2024-11-08', 'AST-020'),
('CMP-047', 'Arista EOS', '4.31.2F', 'Arista Networks', 'Operating System', '2024-06-14', '2024-11-10', 'AST-021'),
('CMP-048', 'Lenovo UEFI BIOS', 'R2IET44W', 'Lenovo', 'UEFI', '2024-06-21', '2024-11-12', 'AST-022'),
('CMP-049', 'FortiOS Firmware', '7.0.14', 'Fortinet', 'Operating System', '2024-06-28', '2024-11-14', 'AST-023'),
('CMP-050', 'HP UEFI BIOS', 'S93 Ver. 01.12.00', 'HP Inc.', 'UEFI', '2024-07-05', NULL, 'AST-024'),
('CMP-051', 'Cisco NX-OS System Software', '10.3(3)', 'Cisco Systems', 'Operating System', '2024-07-12', '2024-11-18', 'AST-025'),
('CMP-052', 'Apple BridgeOS', '22.1.0', 'Apple Inc.', 'Bridge Controller', '2024-07-19', '2024-11-20', 'AST-026'),
('CMP-053', 'HPE UEFI BIOS', 'U54 v1.08', 'HPE', 'UEFI', '2024-07-26', '2024-11-22', 'AST-027'),
('CMP-054', 'Intel ME Firmware', '16.1.30.2307', 'Intel Corporation', 'Firmware', '2024-08-02', '2024-11-24', 'AST-028'),
('CMP-055', 'Juniper QFX Firmware', '20.4R3-S7', 'Juniper Networks', 'Operating System', '2024-08-09', '2024-11-26', 'AST-029'),
('CMP-056', 'Supermicro UEFI BIOS', '2.9', 'Supermicro', 'UEFI', '2024-08-16', '2024-11-28', 'AST-030'),
('CMP-057', 'PAN-OS Threat Prevention', '11.1.2-h3', 'Palo Alto Networks', 'Security Engine', '2024-08-23', '2024-11-30', 'AST-031'),
('CMP-058', 'Cisco IOS-XE System Software', '17.9.4a', 'Cisco Systems', 'Operating System', '2024-08-30', NULL, 'AST-032'),
('CMP-059', 'Lenovo Embedded Controller', '1.30', 'Lenovo', 'Firmware', '2024-09-06', '2024-12-04', 'AST-033'),
('CMP-060', 'Arista Aboot', '8.0.1', 'Arista Networks', 'Boot Firmware', '2024-09-13', '2024-12-06', 'AST-034'),
('CMP-061', 'Dell PowerStore OS', '3.6.0.2', 'Dell Inc.', 'Storage Operating System', '2024-09-20', '2024-12-08', 'AST-035'),
('CMP-062', 'ArubaOS-CX', '10.13.1020', 'HPE Aruba', 'Operating System', '2024-09-27', '2024-12-10', 'AST-036'),
('CMP-063', 'FortiSwitchOS', '7.2.4', 'Fortinet', 'Operating System', '2024-10-04', '2024-12-12', 'AST-037'),
('CMP-064', 'Apple EFI Firmware', '11881.1.2', 'Apple Inc.', 'UEFI', '2024-10-11', NULL, 'AST-038'),
('CMP-065', 'Meraki MX Firmware', '18.211.2', 'Cisco Meraki', 'Operating System', '2024-10-18', '2024-12-16', 'AST-039'),
('CMP-066', 'Lenovo XClarity Controller', '3.30', 'Lenovo', 'Management Controller', '2024-10-25', NULL, 'AST-040')
ON CONFLICT (id) DO NOTHING;

-- Insert Complementary Scans
INSERT INTO scan (id, performedat, scannername, componentid) VALUES
('SCN-039', '2024-10-09', 'Dell Security Suite', 'CMP-039'),
('SCN-040', '2024-10-09', 'PAN-OS Threat Analyzer', 'CMP-040'),
('SCN-041', '2024-10-09', 'XClarity Firmware Inspector', 'CMP-041'),
('SCN-042', '2024-11-02', 'iLO Security Scanner', 'CMP-043'),
('SCN-043', '2024-11-04', 'Junos Security Scanner', 'CMP-044'),
('SCN-044', '2024-11-06', 'BMC Vulnerability Scanner', 'CMP-045'),
('SCN-045', '2024-11-08', 'UEFI Firmware Parser', 'CMP-046'),
('SCN-046', '2024-11-10', 'Arista CloudVision', 'CMP-047'),
('SCN-047', '2024-11-12', 'Lenovo Security Advisor', 'CMP-048'),
('SCN-048', '2024-11-14', 'FortiGuard Scanner', 'CMP-049'),
('SCN-049', '2024-11-18', 'Cisco Security Advisor', 'CMP-051'),
('SCN-050', '2024-11-20', 'Apple Security Scanner', 'CMP-052'),
('SCN-051', '2024-11-22', 'HPE Security Bulletin Scanner', 'CMP-053'),
('SCN-052', '2024-11-24', 'MEAnalyzer', 'CMP-054'),
('SCN-053', '2024-11-26', 'Juniper SIRT Tool', 'CMP-055'),
('SCN-054', '2024-11-28', 'CHIPSEC', 'CMP-056'),
('SCN-055', '2024-11-30', 'PAN-OS Compliance Auditor', 'CMP-057'),
('SCN-056', '2024-12-04', 'Firmware Security Scanner', 'CMP-059'),
('SCN-057', '2024-12-06', 'Aboot Security Checker', 'CMP-060'),
('SCN-058', '2024-12-08', 'Storage Firmware Analyzer', 'CMP-061'),
('SCN-059', '2024-12-10', 'Aruba Security Insights', 'CMP-062'),
('SCN-060', '2024-12-12', 'FortiSwitch Analyzer', 'CMP-063'),
('SCN-061', '2024-12-16', 'Meraki Threat Monitor', 'CMP-065')
ON CONFLICT (id) DO NOTHING;

-- Insert Complementary Vulnerabilities
INSERT INTO vulnerability (id, description, severity, scanid) VALUES
('VUL-041', 'iDRAC9 firmware allows privilege escalation through crafted Redfish payload', 'HIGH', 'SCN-039'),
('VUL-042', 'Junos SRX packet processing flaw allows remote code execution under crafted flow state', 'CRITICAL', 'SCN-043'),
('VUL-043', 'ThinkPad firmware update capsule validation can be bypassed with malformed metadata', 'MEDIUM', 'SCN-047'),
('VUL-044', 'BridgeOS memory handling issue may expose sensitive process data after repeated wake events', 'LOW', 'SCN-050'),
('VUL-045', 'Supermicro UEFI variable protections can be disabled from privileged runtime context', 'HIGH', 'SCN-054'),
('VUL-046', 'Arista boot firmware accepts unsigned recovery image in degraded boot mode', 'CRITICAL', 'SCN-057'),
('VUL-047', 'FortiSwitch management daemon leaks session artifacts in verbose diagnostic mode', 'MEDIUM', 'SCN-060'),
('VUL-048', 'iLO firmware TLS stack susceptible to crafted certificate chain validation bypass', 'HIGH', 'SCN-042'),
('VUL-049', 'Arista EOS control-plane service vulnerable to privilege escalation via malformed API request', 'CRITICAL', 'SCN-046'),
('VUL-050', 'NX-OS image verification process allows rollback to known vulnerable signed build', 'MEDIUM', 'SCN-049'),
('VUL-051', 'QFX firmware permits unauthorized route policy injection through management plane race condition', 'HIGH', 'SCN-053'),
('VUL-052', 'Lenovo embedded controller debug interface remains reachable after secure mode transition', 'LOW', 'SCN-056'),
('VUL-053', 'ArubaOS-CX process isolation weakness allows lateral impact across management services', 'CRITICAL', 'SCN-059')
ON CONFLICT (id) DO NOTHING;

-- Insert Complementary Threats
INSERT INTO threat (id, description, risklevel, type, scanid) VALUES
('THR-041', 'Coordinated exploit campaign against PAN-OS management plane', 'HIGH', 'Remote Exploitation', 'SCN-040'),
('THR-042', 'Compromised IPMI channel enables stealth command-and-control persistence', 'MEDIUM', 'Backdoor', 'SCN-044'),
('THR-043', 'Firewall policy tampering can open unauthorized east-west pathways between segments', 'HIGH', 'Policy Evasion', 'SCN-048'),
('THR-044', 'UEFI misconfiguration exposure increases risk of credential scraping during local access', 'LOW', 'Credential Compromise', 'SCN-051'),
('THR-045', 'Threat actor playbooks target campus edge firewalls for staged data exfiltration', 'HIGH', 'Data Exfiltration', 'SCN-055'),
('THR-046', 'Storage appliance attack chain may enable silent corruption of recovery snapshots', 'MEDIUM', 'Data Integrity', 'SCN-058'),
('THR-047', 'Retail SD-WAN edge compromise can facilitate long-lived unauthorized remote access', 'LOW', 'Remote Access', 'SCN-061'),
('THR-048', 'Exploitation of iLO management endpoints enables persistent out-of-band control', 'HIGH', 'Persistence Threat', 'SCN-042'),
('THR-049', 'Switch control-plane abuse allows selective traffic interception in datacenter spine', 'MEDIUM', 'Man-in-the-Middle', 'SCN-046'),
('THR-050', 'Nexus firmware exploitation may disrupt cluster east-west communication at scale', 'HIGH', 'Network Disruption', 'SCN-049'),
('THR-051', 'Leaf switch compromise enables covert route redirection across storage network', 'LOW', 'Network Hijacking', 'SCN-053'),
('THR-052', 'Endpoint controller takeover can bypass kiosk hardening controls in manufacturing floor', 'HIGH', 'Privilege Escalation', 'SCN-056'),
('THR-053', 'Campus core switch attack chain allows lateral movement across multiple VLAN zones', 'MEDIUM', 'Lateral Movement', 'SCN-059')
ON CONFLICT (id) DO NOTHING;
