DELETE FROM emu_log;
INSERT INTO emu_log
("emu_no", "train_no", "date", "rowid") VALUES

('CR200J2015', 'D5468/D5465', '2020-11-13', 31),
('CR200J2015', 'D5466/D5467', '2020-11-13', 32),
('CR200J2040', 'D5464/1/4', '2020-11-13', 34),
('CR200J2015', 'D5464/1/4', '2020-11-14', 45),
('CR200J2015', 'D5456/D5457', '2020-11-14', 49),
('CR200J2015', 'D5458/D5455', '2020-11-15', 50),
('CR200J2015', 'D5462/D5463', '2020-11-15', 51),
('CR200J2015', 'D5468/D5465', '2020-11-15', 55),
('CR200J2015', 'D5466/D5467', '2020-11-15', 59),
('CR200J2015', 'D5464/1/4', '2020-11-16', 60),
('CR200J2015', 'D5456/D5457', '2020-11-16', 62),
('CR200J2015', 'D5468/D5465', '2020-11-17', 72),

('CRH2A2015', 'D3219', '2020-11-10', 07),
('CRH2A2015', 'D3206', '2020-11-10', 09),
('CRH2A2015', 'D3072/D3073', '2020-11-15', 54),
('CRH2A2015', 'D3074/D3071', '2020-11-15', 57),
('CRH2A2015', 'D3074/D3071', '2020-11-16', 61),
('CRH2A2015', 'D3205', '2020-11-18', 80),
('CRH2A2015', 'D3220', '2020-11-18', 82),

('CR400AF2015', 'G666',  '2020-10-26', 06),
('CR400AF0207', 'G6716', '2020-11-14', 42),
('CR400AF0207', 'G655',  '2020-11-14', 46),
('CR400AF0207', 'G666',  '2020-11-14', 47),
('CR400AF0207', 'G8907', '2020-11-14', 48),
('CR400AF0207', 'G6716', '2020-11-15', 52),
('CR400AF0207', 'G655',  '2020-11-15', 53),
('CR400AF0207', 'G666',  '2020-11-15', 56),
('CR400AF0207', 'G8907', '2020-11-15', 58);

DELETE FROM emu_latest;
INSERT INTO emu_latest
("emu_no", "train_no", "date", "log_id") VALUES
('CR400AF0207', 'G6716', '2020-11-14', 42),
('CR200J2015', 'D5464', '2020-11-14', 45),
('CR200J2015', 'D5461', '2020-11-14', 45),
('CR400AF0207', 'G666',  '2020-11-14', 47),
('CR200J2015', 'D5456', '2020-11-14', 49),
('CR200J2015', 'D5457', '2020-11-14', 49),
('CR200J2015', 'D5458', '2020-11-15', 50),
('CR200J2015', 'D5455', '2020-11-15', 50),
('CR200J2015', 'D5462', '2020-11-15', 51),
('CR200J2015', 'D5463', '2020-11-15', 51),
('CR400AF0207', 'G655',  '2020-11-15', 53),
('CR200J2015', 'D5468', '2020-11-15', 55),
('CR200J2015', 'D5465', '2020-11-15', 55),
('CRH2A2015', 'D3074', '2020-11-15', 57),
('CRH2A2015', 'D3071', '2020-11-15', 57),
('CR400AF0207', 'G8907', '2020-11-15', 58),
('CR200J2015', 'D5466', '2020-11-15', 59),
('CR200J2015', 'D5467', '2020-11-15', 59),
('CRH2A2015', 'D3205', '2020-11-18', 80);

DELETE FROM emu_qrcode;
INSERT INTO emu_qrcode
("emu_no", "emu_bureau", "emu_qrcode", "rowid") VALUES
('CRH2A2015', 'H', 'PQ0504500', 1006),
('CRH2A2015', 'H', 'PQ0505000', 1007),
('CRH2A2015', 'H', 'PQ0558500', 1108),
('CRH2A2015', 'H', 'PQ0731000', 1421),
('CRH2A2015', 'H', 'PQ0731500', 1422),
('CRH2A2015', 'H', 'PQ1375000', 2193),
('CR400AF0207', 'P', '50006000', 2224),
('CR400AF0207', 'P', '50006500', 2225),
('CR400AF2015', 'P', '50009000', 2604),
('CR400AF2015', 'P', '50009500', 2605),
('CR400AF0207', 'P', '50426500', 2671),
('CR200J2015', 'H', 'PQ0916000', 2819),
('CR200J2015', 'H', 'PQ0916500', 2820),
('CR400AF0207', 'P', '50425000', 2969),
('CR400AF2015', 'P', '50370500', 2984),
('CR400AF2015', 'P', '50371000', 2985),
('CR400AF0207', 'P', '60019000', 3363),
('CR400AF0207', 'P', '60120500', 3410),
('CR400AF2015', 'P', '60395000', 3632),
('CR400AF0207', 'P', '60171500', 3833),
('CR400AF0207', 'P', '60442000', 3834),
('CR400AF0207', 'P', '50563500', 3835),
('CH001', 'N', '053', 4764),
('CRH2650', 'N', '111', 4773),
('CRH5A5075', 'N', '472', 4824),
('CR400AF2015', 'P', '60423000', 4925),
('CR400AF0207', 'P', '50704000', 4964),
('CR400AF0207', 'P', '50704500', 4965),
('CR400AF2015', 'P', '50746500', 5005),
('CR200J2040', 'H', 'PQ1630000', 5267),
('CR200J2040', 'H', 'PQ1630500', 5268),
('CR200J2040', 'H', 'PQ1631000', 5269),
('CR400AF2015', 'P', '50880000', 5445);
