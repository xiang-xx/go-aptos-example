insert into "public"."connection" ("chain_a", "token_a", "chain_b", "token_b", "route", "status", "ext")
	values
		('aptos_testnet', 'APTOS_TEST', 'aptos_testnet', 'testBTC', 3, 1, '{"poolInfos":[{"address":"0xe98445b5e7489d1a4afee94940ca4c40e1f6c87a59c3b392e4744614af209de4", "name":"amm"}]}'),
		('aptos_testnet', 'APTOS_TEST', 'aptos_testnet', 'testUSDC', 3, 1, '{"poolInfos":[{"address":"0xe98445b5e7489d1a4afee94940ca4c40e1f6c87a59c3b392e4744614af209de4", "name":"amm"}]}'),
		('aptos_testnet', 'testBTC', 'aptos_testnet', 'testUSDC', 3, 1, '{"poolInfos":[{"address":"0xe98445b5e7489d1a4afee94940ca4c40e1f6c87a59c3b392e4744614af209de4", "name":"amm"}]}');