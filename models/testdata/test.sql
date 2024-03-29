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

DELETE FROM emu_qr_code;
INSERT INTO emu_qr_code
("rowid", "emu_no", "adapter", "operator", "qr_code") VALUES
(12,    'CR400BF3010', 'H', 'H', 'PQ0004200'),
(988,   'CRH380D1585', 'H', 'H', 'PQ0495500'),
(989,   'CRH380D1585', 'H', 'H', 'PQ0496000'),
(1006,  'CRH2A2015',   'H', 'H', 'PQ0504500'),
(1007,  'CRH2A2015',   'H', 'H', 'PQ0505000'),
(1108,  'CRH2A2015',   'H', 'H', 'PQ0558500'),
(1171,  'CR400BF3010', 'H', 'H', 'PQ0591000'),
(1172,  'CR400BF3010', 'H', 'H', 'PQ0591500'),
(1421,  'CRH2A2015',   'H', 'H', 'PQ0731000'),
(1422,  'CRH2A2015',   'H', 'H', 'PQ0731500'),
(1476,  'CRH380D1585', 'H', 'H', 'PQ0764000'),
(2193,  'CRH2A2015',   'H', 'H', 'PQ1375000'),
(2224,  'CR400AF0207', 'P', 'P', '50006000'),
(2225,  'CR400AF0207', 'P', 'P', '50006500'),
(2604,  'CR400AF2015', 'P', 'P', '50009000'),
(2605,  'CR400AF2015', 'P', 'P', '50009500'),
(2671,  'CR400AF0207', 'P', 'P', '50426500'),
(2819,  'CR200J2015',  'H', 'H', 'PQ0916000'),
(2820,  'CR200J2015',  'H', 'H', 'PQ0916500'),
(2969,  'CR400AF0207', 'P', 'P', '50425000'),
(2984,  'CR400AF2015', 'P', 'P', '50370500'),
(2985,  'CR400AF2015', 'P', 'P', '50371000'),
(3363,  'CR400AF0207', 'P', 'P', '60019000'),
(3410,  'CR400AF0207', 'P', 'P', '60120500'),
(3632,  'CR400AF2015', 'P', 'P', '60395000'),
(3766,  'CR400BF3010', 'H', 'H', 'PQ1284500'),
(3833,  'CR400AF0207', 'P', 'P', '60171500'),
(3834,  'CR400AF0207', 'P', 'P', '60442000'),
(3835,  'CR400AF0207', 'P', 'P' ,'50563500'),
(4764,  'CH001',       'F', '?', '053'),
(4773,  'CRH2650',     'F', '?', '111'),
(4824,  'CRH5A5075',   'F', 'N', '472'),
(4925,  'CR400AF2015', 'P', 'P', '60423000'),
(4964,  'CR400AF0207', 'P', 'P', '50704000'),
(4965,  'CR400AF0207', 'P', 'P', '50704500'),
(5005,  'CR400AF2015', 'P', 'P', '50746500'),
(5155,  'CR400BF3010', 'H', 'H', 'PQ1552000'),
(5267,  'CR200J2040',  'H', 'H', 'PQ1630000'),
(5268,  'CR200J2040',  'H', 'H', 'PQ1630500'),
(5269,  'CR200J2040',  'H', 'H', 'PQ1631000'),
(5445,  'CR400AF2015', 'P', 'P', '50880000'),
(5586,  'CR400AF2015', 'P', 'P', '50756500'),
(5993,  'CR300BF3010', 'H', 'H', 'PQ1738000'),
(5994,  'CR300BF3010', 'H', 'H', 'PQ1738500'),
(6753,  'CR300AF2015', 'W', 'W', ',CR300AF,2015,01,,'),
(7023,  'CR200J2015',  'H', 'H', 'PQ2098500'),
(7280,  'CR400AF0207', 'P', 'P', '51066440'),
(7354,  'CR400AF0207', 'P', 'P', '51067000'),
(7415,  'CR400AF0207', 'P', 'P', '51066500'),
(7532,  'CR400AF0207', 'P', 'P', '60807500'),
(8163,  'CR400AF0207', 'P', 'P', '51572500'),
(8174,  'CR400AF0207', 'P', 'P', '51573000'),
(8303,  'CR400AF2015', 'P', 'P', '51737500'),
(8308,  'CR400AF2015', 'P', 'P', '51362000'),
(8309,  'CR400AF2015', 'P', 'P', '51494000'),
(8497,  'CRH3C3010',   'W', 'W', '1069377,CRH3C,3010,01,10,D'),
(8570,  'CR400AF2015', 'P', 'P', '51742500'),
(8656,  'CRH380D1585', 'W', 'W', '1123396,CRH380D,1585,07,08,D'),
(9052,  'CR300BF3010', 'U', 'H', 'PQFB245B95070E4BA9B123D21ED2880C76'),
(9169,  'CRH2A2015',   'U', 'H', 'PV0000079500'),
(9510,  'CR200J2015',  'U', 'H', 'PV0000063800'),
(9511,  'CR200J2015',  'U', 'H', 'PV0000064000'),
(9543,  'CRH3C3010',   'W', 'W', '1069805,CRH3C,3010,07,13,A'),
(10066, 'CR400AF0207', 'M', 'P', 'CR400AF-0207-05-03A'),
(10309, 'CR400BF3010', 'U', 'H', 'PV0000104100'),
(10440, 'CRH2A2015',   'U', 'H', 'PV0000180300'),
(10569, 'CR400BF3010', 'U', 'H', 'PV0000197900'),
(10570, 'CR400BF3010', 'U', 'H', 'PV0000198000'),
(10687, 'CR400AF2015', 'M', 'P', 'CR400AF-2015-03-10A'),
(10834, 'CR300AF2015', 'M', 'W', 'CR300AF-2015-06-08A'),
(10940, 'CR300BF3010', 'M', 'H', 'H0140000'),
(11023, 'CR400BF3010', 'M', 'H', 'H0035000'),
(11030, 'CRH380D1585', 'M', 'W', 'W0006000'),
(11073, 'CRH3C3010',   'M', 'W', 'W0144000');
