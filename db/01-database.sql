-- Create Severity Lookup Table
CREATE TABLE IF NOT EXISTS severity (
    severitylevel VARCHAR(20) PRIMARY KEY
);

-- Create RiskLevel Lookup Table
CREATE TABLE IF NOT EXISTS risklevel (
    risklevel VARCHAR(20) PRIMARY KEY
);

-- Insert Severity Levels
INSERT INTO severity (severitylevel) VALUES 
('LOW'), 
('MEDIUM'), 
('HIGH'), 
('CRITICAL')
ON CONFLICT (severitylevel) DO NOTHING;

-- Insert Risk Levels
INSERT INTO risklevel (risklevel) VALUES 
('LOW'), 
('MEDIUM'), 
('HIGH')
ON CONFLICT (risklevel) DO NOTHING;

-- Create Asset Table
CREATE TABLE IF NOT EXISTS asset (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    createdat DATE NOT NULL,
    lastscan DATE
);

-- Create Component Table
CREATE TABLE IF NOT EXISTS component (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(100) NOT NULL,
    vendor VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    createdat DATE NOT NULL,
    lastscan DATE,
    assetid VARCHAR(255) NOT NULL,
    FOREIGN KEY (assetid) REFERENCES asset(id)
);

-- Create Scan Table
CREATE TABLE IF NOT EXISTS scan (
    id VARCHAR(255) PRIMARY KEY,
    performedat DATE NOT NULL,
    scannername VARCHAR(255) NOT NULL,
    componentid VARCHAR(255) NOT NULL,
    FOREIGN KEY (componentid) REFERENCES component(id)
);

-- Create Vulnerability Table
CREATE TABLE IF NOT EXISTS vulnerability (
    id VARCHAR(255) PRIMARY KEY,
    description TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL,
    scanid VARCHAR(255) NOT NULL,
    FOREIGN KEY (severity) REFERENCES severity(severitylevel),
    FOREIGN KEY (scanid) REFERENCES scan(id)
);

-- Create Threat Table
CREATE TABLE IF NOT EXISTS threat (
    id VARCHAR(255) PRIMARY KEY,
    description TEXT NOT NULL,
    risklevel VARCHAR(20) NOT NULL,
    type VARCHAR(100) NOT NULL,
    scanid VARCHAR(255) NOT NULL,
    FOREIGN KEY (risklevel) REFERENCES risklevel(risklevel),
    FOREIGN KEY (scanid) REFERENCES scan(id)
);

-- Add index for performance optimization
CREATE INDEX IF NOT EXISTS idx_component_asset ON component(assetid);
CREATE INDEX IF NOT EXISTS idx_scan_component ON scan(componentid);
CREATE INDEX IF NOT EXISTS idx_scan_component_performedat_id ON scan(componentid, performedat DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_vulnerability_scan ON vulnerability(scanid);
CREATE INDEX IF NOT EXISTS idx_threat_scan ON threat(scanid);
