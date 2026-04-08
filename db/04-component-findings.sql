-- New Dataset for Component and Findings filtering/sorting UI scenarios
-- Goal: 1 Asset with multiple components, each component with multiple vulnerabilities and threats.

-- Insert New Asset
INSERT INTO asset (id, name, description, createdat, lastscan) VALUES
('AST-041', 'HPE Nimble Storage Cluster HF20', 'All-flash adaptive storage array for VDI and transactional workloads', '2024-11-01', '2024-12-20')
ON CONFLICT (id) DO NOTHING;

-- Insert Components for AST-041
INSERT INTO component (id, name, version, vendor, type, createdat, lastscan, assetid) VALUES
('CMP-067', 'Nimble OS Controller Firmware', '6.1.1.0', 'HPE Nimble Storage', 'Storage OS', '2024-11-01', '2024-12-20', 'AST-041'),
('CMP-068', 'Intel X550-T2 NIC Controller', '21.5.0', 'Intel Corporation', 'Network Adapter', '2024-11-01', '2024-12-20', 'AST-041'),
('CMP-069', 'Samsung PM1733 Enterprise NVMe', 'GPK1601Q', 'Samsung Electronics', 'Storage Device', '2024-11-01', '2024-12-20', 'AST-041')
ON CONFLICT (id) DO NOTHING;

-- Insert Scans for AST-041 Components
INSERT INTO scan (id, performedat, scannername, componentid) VALUES
('SCN-062', '2024-12-20', 'Nimble Security Analyzer', 'CMP-067'),
('SCN-063', '2024-12-20', 'Intel Firmware Inspector', 'CMP-068'),
('SCN-064', '2024-12-20', 'NVMe Security Tool', 'CMP-069')
ON CONFLICT (id) DO NOTHING;

-- Insert Vulnerabilities for SCN-062 (CMP-067)
INSERT INTO vulnerability (id, description, severity, scanid) VALUES
('VUL-054', 'Nimble OS allows unauthorized SSH access via deprecated key exchange', 'HIGH', 'SCN-062'),
('VUL-055', 'Buffer overflow in management API allows denial of service', 'MEDIUM', 'SCN-062'),
('VUL-056', 'Information disclosure through verbose SNMP trap messages', 'LOW', 'SCN-062'),
('VUL-057', 'Cross-site scripting (XSS) in local web-based administration console', 'MEDIUM', 'SCN-062')
ON CONFLICT (id) DO NOTHING;

-- Insert Vulnerabilities for SCN-063 (CMP-068)
INSERT INTO vulnerability (id, description, severity, scanid) VALUES
('VUL-058', 'Intel NIC firmware susceptible to side-channel memory timing analysis', 'HIGH', 'SCN-063'),
('VUL-059', 'Malformed Ethernet frame may cause kernel panic via DMA exception', 'CRITICAL', 'SCN-063'),
('VUL-060', 'Unsigned firmware update capsule accepted in manufacturing mode', 'HIGH', 'SCN-063')
ON CONFLICT (id) DO NOTHING;

-- Insert Vulnerabilities for SCN-064 (CMP-069)
INSERT INTO vulnerability (id, description, severity, scanid) VALUES
('VUL-061', 'NVMe SSD controller contains unmapped internal debug registers', 'MEDIUM', 'SCN-064'),
('VUL-062', 'Wear-leveling metadata exposure through vendor-specific commands', 'LOW', 'SCN-064'),
('VUL-063', 'Firmware integrity check bypass using crafted secure erase payload', 'CRITICAL', 'SCN-064')
ON CONFLICT (id) DO NOTHING;

-- Insert Threats for SCN-062 (CMP-067)
INSERT INTO threat (id, description, risklevel, type, scanid) VALUES
('THR-054', 'Ransomware campaign targeting Nimble storage snapshots', 'HIGH', 'Ransomware', 'SCN-062'),
('THR-055', 'Compromised admin credentials via brute-force of management plane', 'MEDIUM', 'Credential Compromise', 'SCN-062'),
('THR-056', 'Exploitation of management API for lateral movement in storage fabric', 'HIGH', 'Lateral Movement', 'SCN-062')
ON CONFLICT (id) DO NOTHING;

-- Insert Threats for SCN-063 (CMP-068)
INSERT INTO threat (id, description, risklevel, type, scanid) VALUES
('THR-057', 'Man-in-the-middle attack via rogue frame injection at NIC layer', 'MEDIUM', 'Man-in-the-Middle', 'SCN-063'),
('THR-058', 'Network interface takeover for stealthy data exfiltration', 'HIGH', 'Data Exfiltration', 'SCN-063'),
('THR-059', 'Persistent firmware backdoor surviving NIC hardware reset', 'HIGH', 'Persistence Threat', 'SCN-063')
ON CONFLICT (id) DO NOTHING;

-- Insert Threats for SCN-064 (CMP-069)
INSERT INTO threat (id, description, risklevel, type, scanid) VALUES
('THR-060', 'Data integrity attack targeting secure erase implementation', 'MEDIUM', 'Data Tampering', 'SCN-064'),
('THR-061', 'Nation-state actors exploiting SSD debug ports for out-of-band access', 'HIGH', 'Advanced Persistent Threat', 'SCN-064'),
('THR-062', 'Supply chain implant in SSD controller firmware logic', 'HIGH', 'Supply Chain Attack', 'SCN-064')
ON CONFLICT (id) DO NOTHING;
