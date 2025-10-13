-- Create Severity Lookup Table
CREATE TABLE severity (
    severitylevel VARCHAR(20) PRIMARY KEY
);

-- Create RiskLevel Lookup Table
CREATE TABLE risklevel (
    risklevel VARCHAR(20) PRIMARY KEY
);

-- Insert Severity Levels
INSERT INTO severity (severitylevel) VALUES 
('LOW'), 
('MEDIUM'), 
('HIGH'), 
('CRITICAL');

-- Insert Risk Levels
INSERT INTO risklevel (risklevel) VALUES 
('LOW'), 
('MEDIUM'), 
('HIGH');

-- Create Asset Table
CREATE TABLE asset (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    createdat DATE,
    lastscan DATE
);

-- Create Component Table
CREATE TABLE component (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    version VARCHAR(100),
    vendor VARCHAR(255),
    type VARCHAR(100),
    createdat DATE,
    lastscan DATE,
    assetid VARCHAR(255),
    FOREIGN KEY (assetid) REFERENCES asset(id)
);

-- Create Scan Table
CREATE TABLE scan (
    id VARCHAR(255) PRIMARY KEY,
    performedat DATE,
    scannername VARCHAR(255),
    componentid VARCHAR(255),
    FOREIGN KEY (componentid) REFERENCES component(id)
);

-- Create Vulnerability Table
CREATE TABLE vulnerability (
    id VARCHAR(255) PRIMARY KEY,
    description TEXT,
    severity VARCHAR(20),
    scanid VARCHAR(255),
    FOREIGN KEY (severity) REFERENCES severity(severitylevel),
    FOREIGN KEY (scanid) REFERENCES scan(id)
);

-- Create Threat Table
CREATE TABLE threat (
    id VARCHAR(255) PRIMARY KEY,
    description TEXT,
    risklevel VARCHAR(20),
    type VARCHAR(100),
    scanid VARCHAR(255),
    FOREIGN KEY (risklevel) REFERENCES risklevel(risklevel),
    FOREIGN KEY (scanid) REFERENCES scan(id)
);

-- Add index for performance optimization
CREATE INDEX idx_component_asset ON component(assetid);
CREATE INDEX idx_scan_component ON scan(componentid);
CREATE INDEX idx_vulnerability_scan ON vulnerability(scanid);
CREATE INDEX idx_threat_scan ON threat(scanid);
