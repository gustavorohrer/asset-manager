-- Sample Data for Hardware Supply Chain Security Database

-- Insert Assets (Hardware Devices)
INSERT INTO asset (id, name, description, createdat, lastscan) VALUES
('AST-001', 'Dell PowerEdge R740 Server', 'Production database server in datacenter rack A3', '2024-01-15', '2024-10-08'),
('AST-002', 'HP EliteBook 850 G8', 'Executive laptop assigned to Finance Department', '2024-02-10', '2024-10-07'),
('AST-003', 'Cisco Catalyst 9300 Switch', 'Core network switch for building 2 distribution layer', '2024-01-20', '2024-10-08'),
('AST-004', 'Lenovo ThinkPad X1 Carbon Gen 9', 'Engineering workstation for software development team', '2024-03-05', '2024-10-06'),
('AST-005', 'Supermicro SYS-2029U Server', 'Virtualization host server in cloud infrastructure', '2024-01-25', '2024-10-08'),
('AST-006', 'Juniper EX4300 Switch', 'Access layer switch for office floor 3', '2024-02-15', '2024-10-07'),
('AST-007', 'Arista 7050X3 Router', 'Border gateway router for WAN connectivity', '2024-01-30', '2024-10-08'),
('AST-008', 'Dell Precision 5820 Workstation', 'CAD workstation for design team', '2024-03-20', '2024-10-05'),
('AST-009', 'Cisco ASR 1001-X Router', 'Backup WAN router for redundancy', '2024-02-05', '2024-10-07'),
('AST-010', 'HPE ProLiant DL380 Gen10', 'Web application server in DMZ', '2024-01-18', '2024-10-08'),
('AST-011', 'MacBook Pro 16-inch M1', 'Creative team laptop for graphic design', '2024-04-01', '2024-10-06'),
('AST-012', 'Fortinet FortiGate 600E', 'Next-generation firewall at network perimeter', '2024-01-22', '2024-10-08');

-- Insert Components (Firmware and Hardware Components)
INSERT INTO component (id, name, version, vendor, type, createdat, lastscan, assetid) VALUES
-- Dell PowerEdge R740 Components
('CMP-001', 'Dell UEFI BIOS', '2.10.2', 'Dell Inc.', 'UEFI', '2024-01-15', '2024-10-08', 'AST-001'),
('CMP-002', 'Intel Management Engine', '4.1.4.50', 'Intel Corporation', 'Firmware', '2024-01-15', '2024-10-08', 'AST-001'),
('CMP-003', 'PERC H740P RAID Controller', '25.5.7.0005', 'Broadcom/LSI', 'Storage Controller', '2024-01-15', '2024-10-08', 'AST-001'),
('CMP-004', 'iDRAC9 BMC', '4.40.00.00', 'Dell Inc.', 'BMC', '2024-01-15', '2024-10-08', 'AST-001'),
('CMP-005', 'Intel X710 NIC Firmware', '7.10', 'Intel Corporation', 'Network Adapter', '2024-01-15', '2024-10-08', 'AST-001'),

-- HP EliteBook 850 G8 Components
('CMP-006', 'HP UEFI BIOS', 'T77 Ver. 01.08.00', 'HP Inc.', 'UEFI', '2024-02-10', '2024-10-07', 'AST-002'),
('CMP-007', 'Intel ME Firmware', '15.0.22.1586', 'Intel Corporation', 'Firmware', '2024-02-10', '2024-10-07', 'AST-002'),
('CMP-008', 'Samsung NVMe SSD Firmware', '5L2QFXV7', 'Samsung Electronics', 'Storage Device', '2024-02-10', '2024-10-07', 'AST-002'),
('CMP-009', 'Intel AX201 WiFi Firmware', '22.100.0.2', 'Intel Corporation', 'Wireless Adapter', '2024-02-10', '2024-10-07', 'AST-002'),
('CMP-010', 'TPM 2.0 Firmware', '7.2.3.0', 'Infineon', 'Security Chip', '2024-02-10', '2024-10-07', 'AST-002'),

-- Cisco Catalyst 9300 Components
('CMP-011', 'IOS-XE Boot Loader', '17.3.4', 'Cisco Systems', 'Boot Firmware', '2024-01-20', '2024-10-08', 'AST-003'),
('CMP-012', 'IOS-XE System Software', '17.3.4a', 'Cisco Systems', 'Operating System', '2024-01-20', '2024-10-08', 'AST-003'),
('CMP-013', 'Cisco StackWise-480 Firmware', '1.2.5', 'Cisco Systems', 'Stacking Module', '2024-01-20', '2024-10-08', 'AST-003'),

-- Lenovo ThinkPad X1 Carbon Components
('CMP-014', 'Lenovo UEFI BIOS', 'N32ET71W', 'Lenovo', 'UEFI', '2024-03-05', '2024-10-06', 'AST-004'),
('CMP-015', 'Intel ME Firmware', '14.1.53.1649', 'Intel Corporation', 'Firmware', '2024-03-05', '2024-10-06', 'AST-004'),
('CMP-016', 'Realtek PCIe GbE Controller', '10.045.1003.2020', 'Realtek', 'Network Adapter', '2024-03-05', '2024-10-06', 'AST-004'),
('CMP-017', 'Intel UHD Graphics Firmware', '27.20.100.9316', 'Intel Corporation', 'GPU', '2024-03-05', '2024-10-06', 'AST-004'),

-- Supermicro Server Components
('CMP-018', 'Supermicro UEFI BIOS', '3.3', 'Supermicro', 'UEFI', '2024-01-25', '2024-10-08', 'AST-005'),
('CMP-019', 'IPMI BMC Firmware', '01.73.10', 'Supermicro', 'BMC', '2024-01-25', '2024-10-08', 'AST-005'),
('CMP-020', 'Broadcom 57414 NIC', '214.0.215.0', 'Broadcom', 'Network Adapter', '2024-01-25', '2024-10-08', 'AST-005'),
('CMP-021', 'LSI MegaRAID SAS-3 3108', '4.680.00-8321', 'Broadcom/LSI', 'RAID Controller', '2024-01-25', '2024-10-08', 'AST-005'),

-- Juniper Switch Components
('CMP-022', 'Junos OS', '18.4R2-S3', 'Juniper Networks', 'Operating System', '2024-02-15', '2024-10-07', 'AST-006'),
('CMP-023', 'Juniper Boot Loader', '12.0X47-D15.4', 'Juniper Networks', 'Boot Firmware', '2024-02-15', '2024-10-07', 'AST-006'),

-- Arista Router Components
('CMP-024', 'EOS Operating System', '4.25.4M', 'Arista Networks', 'Operating System', '2024-01-30', '2024-10-08', 'AST-007'),
('CMP-025', 'Aboot Boot Loader', '6.0.9', 'Arista Networks', 'Boot Firmware', '2024-01-30', '2024-10-08', 'AST-007'),

-- Dell Precision Workstation Components
('CMP-026', 'Dell UEFI BIOS', '2.8.0', 'Dell Inc.', 'UEFI', '2024-03-20', '2024-10-05', 'AST-008'),
('CMP-027', 'NVIDIA Quadro RTX 4000 VBIOS', '90.04.4B.00.48', 'NVIDIA', 'GPU', '2024-03-20', '2024-10-05', 'AST-008'),
('CMP-028', 'Intel Xeon C600 PCU', '6.01', 'Intel Corporation', 'Power Management', '2024-03-20', '2024-10-05', 'AST-008'),

-- Cisco ASR Router Components
('CMP-029', 'IOS-XE System Software', '16.12.4', 'Cisco Systems', 'Operating System', '2024-02-05', '2024-10-07', 'AST-009'),
('CMP-030', 'ROMMON Boot Loader', '16.10(1r)', 'Cisco Systems', 'Boot Firmware', '2024-02-05', '2024-10-07', 'AST-009'),

-- HPE ProLiant Components
('CMP-031', 'HPE UEFI BIOS', 'U30 v2.54', 'HPE', 'UEFI', '2024-01-18', '2024-10-08', 'AST-010'),
('CMP-032', 'iLO 5 Firmware', '2.44', 'HPE', 'BMC', '2024-01-18', '2024-10-08', 'AST-010'),
('CMP-033', 'HPE Smart Array P408i', '5.00', 'HPE', 'RAID Controller', '2024-01-18', '2024-10-08', 'AST-010'),

-- MacBook Pro Components
('CMP-034', 'Apple T2 Security Chip', '19.16.10.0.0', 'Apple Inc.', 'Security Chip', '2024-04-01', '2024-10-06', 'AST-011'),
('CMP-035', 'iBridge Firmware', '19.16.10077.0.0', 'Apple Inc.', 'Bridge Controller', '2024-04-01', '2024-10-06', 'AST-011'),
('CMP-036', 'Apple SSD Controller', '1711.81.1', 'Apple Inc.', 'Storage Controller', '2024-04-01', '2024-10-06', 'AST-011'),

-- Fortinet Firewall Components
('CMP-037', 'FortiOS Firmware', '6.4.7', 'Fortinet', 'Operating System', '2024-01-22', '2024-10-08', 'AST-012'),
('CMP-038', 'FortiASIC NP6 Firmware', '6.4.0', 'Fortinet', 'ASIC Processor', '2024-01-22', '2024-10-08', 'AST-012');

-- Insert Scans
INSERT INTO scan (id, performedat, scannername, componentid) VALUES
('SCN-001', '2024-10-08', 'CHIPSEC', 'CMP-001'),
('SCN-002', '2024-10-08', 'Binwalk', 'CMP-002'),
('SCN-003', '2024-10-08', 'Firmware Analysis Toolkit', 'CMP-003'),
('SCN-004', '2024-10-08', 'FwAnalyzer', 'CMP-004'),
('SCN-005', '2024-10-08', 'CHIPSEC', 'CMP-005'),
('SCN-006', '2024-10-07', 'HP Image Assistant', 'CMP-006'),
('SCN-007', '2024-10-07', 'MEAnalyzer', 'CMP-007'),
('SCN-008', '2024-10-07', 'NVMe CLI Scanner', 'CMP-008'),
('SCN-009', '2024-10-07', 'Firmware Security Scanner', 'CMP-009'),
('SCN-010', '2024-10-07', 'TPM Toolkit', 'CMP-010'),
('SCN-011', '2024-10-08', 'Cisco Security Advisor', 'CMP-011'),
('SCN-012', '2024-10-08', 'IOS-XE Vulnerability Scanner', 'CMP-012'),
('SCN-013', '2024-10-08', 'Cisco PSIRT Scanner', 'CMP-013'),
('SCN-014', '2024-10-06', 'Lenovo Security Advisor', 'CMP-014'),
('SCN-015', '2024-10-06', 'Intel SA-00086 Scanner', 'CMP-015'),
('SCN-016', '2024-10-06', 'PCI Device Scanner', 'CMP-016'),
('SCN-017', '2024-10-06', 'GPU Firmware Scanner', 'CMP-017'),
('SCN-018', '2024-10-08', 'UEFI Firmware Parser', 'CMP-018'),
('SCN-019', '2024-10-08', 'BMC Vulnerability Scanner', 'CMP-019'),
('SCN-020', '2024-10-08', 'Broadcom Security Scanner', 'CMP-020'),
('SCN-021', '2024-10-08', 'RAID Controller Analyzer', 'CMP-021'),
('SCN-022', '2024-10-07', 'Junos Security Scanner', 'CMP-022'),
('SCN-023', '2024-10-07', 'Juniper SIRT Tool', 'CMP-023'),
('SCN-024', '2024-10-08', 'Arista CloudVision', 'CMP-024'),
('SCN-025', '2024-10-08', 'Aboot Security Checker', 'CMP-025'),
('SCN-026', '2024-10-05', 'CHIPSEC', 'CMP-026'),
('SCN-027', '2024-10-05', 'GPU-Z Validator', 'CMP-027'),
('SCN-028', '2024-10-05', 'Intel PCU Scanner', 'CMP-028'),
('SCN-029', '2024-10-07', 'Cisco IOS Scanner', 'CMP-029'),
('SCN-030', '2024-10-07', 'ROMMON Validator', 'CMP-030'),
('SCN-031', '2024-10-08', 'HPE Security Bulletin Scanner', 'CMP-031'),
('SCN-032', '2024-10-08', 'iLO Security Scanner', 'CMP-032'),
('SCN-033', '2024-10-08', 'Smart Array Firmware Tool', 'CMP-033'),
('SCN-034', '2024-10-06', 'Apple Security Scanner', 'CMP-034'),
('SCN-035', '2024-10-06', 'iBridge Analyzer', 'CMP-035'),
('SCN-036', '2024-10-06', 'SSD Health Monitor', 'CMP-036'),
('SCN-037', '2024-10-08', 'FortiGuard Scanner', 'CMP-037'),
('SCN-038', '2024-10-08', 'FortiASIC Validator', 'CMP-038');

-- Insert Vulnerabilities
INSERT INTO vulnerability (id, description, severity, scanid) VALUES
('VUL-001', 'Dell UEFI BIOS vulnerable to BootHole (CVE-2020-10713) allowing SecureBoot bypass', 'CRITICAL', 'SCN-001'),
('VUL-002', 'Intel Management Engine contains critical privilege escalation flaw (CVE-2020-8705)', 'CRITICAL', 'SCN-002'),
('VUL-003', 'PERC RAID controller firmware has unpatched buffer overflow vulnerability', 'HIGH', 'SCN-003'),
('VUL-004', 'iDRAC9 BMC susceptible to authentication bypass (CVE-2021-21514)', 'CRITICAL', 'SCN-004'),
('VUL-005', 'Intel X710 NIC firmware contains DMA attack vulnerability', 'HIGH', 'SCN-005'),
('VUL-006', 'HP UEFI BIOS missing patches for multiple CVEs including memory corruption', 'HIGH', 'SCN-006'),
('VUL-007', 'Intel ME firmware version vulnerable to INTEL-SA-00213 privilege escalation', 'CRITICAL', 'SCN-007'),
('VUL-008', 'Samsung SSD firmware susceptible to data corruption under power loss', 'MEDIUM', 'SCN-008'),
('VUL-009', 'Intel WiFi firmware allows packet injection attacks', 'MEDIUM', 'SCN-009'),
('VUL-010', 'TPM 2.0 implementation vulnerable to timing attack TPM-FAIL', 'HIGH', 'SCN-010'),
('VUL-011', 'Cisco IOS-XE contains hardcoded credentials in ROMMON (CVE-2021-1435)', 'CRITICAL', 'SCN-011'),
('VUL-012', 'IOS-XE system software vulnerable to remote code execution (CVE-2021-34770)', 'CRITICAL', 'SCN-012'),
('VUL-013', 'StackWise firmware allows unauthorized stack member addition', 'MEDIUM', 'SCN-013'),
('VUL-014', 'Lenovo UEFI BIOS vulnerable to SMM privilege escalation', 'HIGH', 'SCN-014'),
('VUL-015', 'Intel ME firmware contains INTEL-SA-00213 critical vulnerability', 'CRITICAL', 'SCN-015'),
('VUL-016', 'Realtek NIC firmware has exploitable buffer overflow', 'MEDIUM', 'SCN-016'),
('VUL-017', 'Intel UHD Graphics firmware vulnerable to information disclosure', 'LOW', 'SCN-017'),
('VUL-018', 'Supermicro UEFI BIOS allows SecureBoot policy manipulation', 'CRITICAL', 'SCN-018'),
('VUL-019', 'IPMI BMC firmware vulnerable to cipher zero authentication bypass', 'CRITICAL', 'SCN-019'),
('VUL-020', 'Broadcom NIC firmware contains remote code execution vulnerability', 'HIGH', 'SCN-020'),
('VUL-021', 'LSI MegaRAID controller susceptible to firmware backdoor implantation', 'HIGH', 'SCN-021'),
('VUL-022', 'Junos OS vulnerable to BGP hijacking attack (CVE-2020-1631)', 'HIGH', 'SCN-022'),
('VUL-023', 'Juniper boot loader allows unsigned code execution', 'CRITICAL', 'SCN-023'),
('VUL-024', 'Arista EOS vulnerable to privilege escalation via CLI', 'HIGH', 'SCN-024'),
('VUL-025', 'Aboot allows unsigned kernel loading bypassing secure boot', 'CRITICAL', 'SCN-025'),
('VUL-026', 'Dell Precision BIOS vulnerable to SmmSmramSaveState attack', 'HIGH', 'SCN-026'),
('VUL-027', 'NVIDIA GPU VBIOS contains exploitable firmware parsing vulnerability', 'MEDIUM', 'SCN-027'),
('VUL-028', 'Intel Xeon PCU firmware has unpatched voltage manipulation flaw', 'MEDIUM', 'SCN-028'),
('VUL-029', 'Cisco IOS-XE RESTCONF API authentication bypass (CVE-2021-1619)', 'CRITICAL', 'SCN-029'),
('VUL-030', 'ROMMON contains hardcoded encryption keys', 'HIGH', 'SCN-030'),
('VUL-031', 'HPE UEFI BIOS vulnerable to BootHole and related vulnerabilities', 'CRITICAL', 'SCN-031'),
('VUL-032', 'iLO 5 firmware susceptible to remote code execution (CVE-2020-7200)', 'CRITICAL', 'SCN-032'),
('VUL-033', 'Smart Array controller firmware allows unauthorized disk access', 'HIGH', 'SCN-033'),
('VUL-034', 'Apple T2 chip vulnerable to checkm8 exploit allowing jailbreak', 'HIGH', 'SCN-034'),
('VUL-035', 'iBridge firmware contains race condition in DMA protection', 'MEDIUM', 'SCN-035'),
('VUL-036', 'Apple SSD controller firmware has wear-leveling vulnerability', 'LOW', 'SCN-036'),
('VUL-037', 'FortiOS firmware vulnerable to SSL VPN pre-authentication RCE (CVE-2022-42475)', 'CRITICAL', 'SCN-037'),
('VUL-038', 'FortiASIC firmware contains buffer overflow in packet processing', 'HIGH', 'SCN-038'),
('VUL-039', 'Dell iDRAC allows unauthorized firmware updates without authentication', 'CRITICAL', 'SCN-004'),
('VUL-040', 'Supermicro BMC contains default credentials not disabled in production', 'HIGH', 'SCN-019');

-- Insert Threats
INSERT INTO threat (id, description, risklevel, type, scanid) VALUES
('THR-001', 'Bootkits and rootkits can bypass SecureBoot via BootHole exploit', 'HIGH', 'Firmware Implant', 'SCN-001'),
('THR-002', 'Nation-state actors exploiting ME vulnerability for persistent access', 'HIGH', 'Advanced Persistent Threat', 'SCN-002'),
('THR-003', 'Ransomware targeting RAID controller to encrypt storage arrays', 'HIGH', 'Ransomware', 'SCN-003'),
('THR-004', 'Remote attackers gaining full system control via BMC exploitation', 'HIGH', 'Remote Access', 'SCN-004'),
('THR-005', 'DMA attacks allowing memory dumping and credential theft', 'HIGH', 'Hardware Attack', 'SCN-005'),
('THR-006', 'Supply chain attack via compromised BIOS update mechanism', 'HIGH', 'Supply Chain', 'SCN-006'),
('THR-007', 'Firmware-level implants surviving OS reinstallation via ME compromise', 'HIGH', 'Persistence Threat', 'SCN-007'),
('THR-008', 'Data loss and corruption through malicious SSD firmware manipulation', 'MEDIUM', 'Data Integrity', 'SCN-008'),
('THR-009', 'Man-in-the-middle attacks via rogue WiFi firmware injection', 'MEDIUM', 'Network Attack', 'SCN-009'),
('THR-010', 'Cryptographic key extraction via TPM timing side-channel', 'MEDIUM', 'Cryptographic Attack', 'SCN-010'),
('THR-011', 'Network infrastructure takeover via hardcoded Cisco credentials', 'HIGH', 'Credential Compromise', 'SCN-011'),
('THR-012', 'Remote exploitation of IOS-XE enabling traffic interception', 'HIGH', 'Network Compromise', 'SCN-012'),
('THR-013', 'Unauthorized network segmentation breach via stack manipulation', 'MEDIUM', 'Network Attack', 'SCN-013'),
('THR-014', 'Privilege escalation to SMM allowing complete system compromise', 'HIGH', 'Privilege Escalation', 'SCN-014'),
('THR-015', 'Nation-state surveillance via Intel ME persistent backdoor', 'HIGH', 'Espionage', 'SCN-015'),
('THR-016', 'Network card compromise enabling stealth data exfiltration', 'MEDIUM', 'Data Exfiltration', 'SCN-016'),
('THR-017', 'GPU firmware exploitation for cryptocurrency mining malware', 'LOW', 'Cryptojacking', 'SCN-017'),
('THR-018', 'Pre-boot malware installation bypassing all OS-level security', 'HIGH', 'Bootkit', 'SCN-018'),
('THR-019', 'BMC exploitation for out-of-band persistent remote access', 'HIGH', 'Backdoor', 'SCN-019'),
('THR-020', 'Network adapter firmware implant for traffic interception', 'HIGH', 'Man-in-the-Middle', 'SCN-020'),
('THR-021', 'RAID controller backdoor enabling covert data modification', 'HIGH', 'Data Tampering', 'SCN-021'),
('THR-022', 'BGP route hijacking leading to traffic redirection', 'HIGH', 'Network Hijacking', 'SCN-022'),
('THR-023', 'Malicious boot loader replacing legitimate OS kernel', 'HIGH', 'Boot Attack', 'SCN-023'),
('THR-024', 'CLI exploitation for lateral movement within network infrastructure', 'MEDIUM', 'Lateral Movement', 'SCN-024'),
('THR-025', 'Secure boot bypass enabling persistent firmware malware', 'HIGH', 'Security Bypass', 'SCN-025'),
('THR-026', 'SMM-based rootkit surviving firmware updates and OS changes', 'HIGH', 'Advanced Rootkit', 'SCN-026'),
('THR-027', 'GPU firmware manipulation for display output interception', 'LOW', 'Privacy Threat', 'SCN-027'),
('THR-028', 'Voltage fault injection attacks via PCU manipulation', 'MEDIUM', 'Physical Attack', 'SCN-028'),
('THR-029', 'Zero-day exploitation of RESTCONF API for router compromise', 'HIGH', 'Zero-Day Exploit', 'SCN-029'),
('THR-030', 'Hardcoded key extraction enabling encrypted traffic decryption', 'HIGH', 'Encryption Defeat', 'SCN-030'),
('THR-031', 'Multi-stage bootkit chain compromising entire boot process', 'HIGH', 'Boot Chain Attack', 'SCN-031'),
('THR-032', 'Out-of-band management exploitation for stealthy persistence', 'HIGH', 'Stealth Attack', 'SCN-032'),
('THR-033', 'Storage controller compromise enabling data theft without detection', 'HIGH', 'Data Breach', 'SCN-033'),
('THR-034', 'T2 chip exploit allowing macOS security feature bypass', 'MEDIUM', 'Platform Security Bypass', 'SCN-034'),
('THR-035', 'DMA protection bypass enabling direct memory access attacks', 'MEDIUM', 'Memory Attack', 'SCN-035'),
('THR-036', 'SSD wear-leveling exploitation for deleted data recovery', 'LOW', 'Data Recovery Threat', 'SCN-036'),
('THR-037', 'SSL VPN exploitation for unauthorized network access', 'HIGH', 'Network Intrusion', 'SCN-037'),
('THR-038', 'Hardware accelerator compromise affecting all encrypted traffic', 'HIGH', 'Encryption Compromise', 'SCN-038'),
('THR-039', 'Supply chain firmware injection during manufacturing process', 'HIGH', 'Supply Chain Attack', 'SCN-001'),
('THR-040', 'Default credential exploitation for BMC remote control', 'MEDIUM', 'Weak Authentication', 'SCN-019');