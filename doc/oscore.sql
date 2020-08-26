DROP TABLE IF EXISTS `tbl_api_crypto_info`;
DROP TABLE IF EXISTS `tbl_order_api_comment`;
DROP TABLE IF EXISTS `tbl_tool_box`;
DROP TABLE IF EXISTS `tbl_algorithm_env`;
DROP TABLE IF EXISTS `tbl_api_algorithm`;
DROP TABLE IF EXISTS `tbl_env`;
DROP TABLE IF EXISTS `tbl_algorithm`;
DROP TABLE IF EXISTS `tbl_api_test_key`;
DROP TABLE IF EXISTS `tbl_qr_code`;
DROP TABLE IF EXISTS `tbl_api_key`;
DROP TABLE IF EXISTS `tbl_order`;
DROP TABLE IF EXISTS `tbl_error_code`;
DROP TABLE IF EXISTS `tbl_request_param`;
DROP TABLE IF EXISTS `tbl_specifications`;
DROP TABLE IF EXISTS `tbl_api_basic_info`;
DROP TABLE IF EXISTS `tbl_score`;

create table tbl_api_basic_info
(
 ApiId INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 ApiType varchar(255) NOT NULL DEFAULT '' COMMENT '',
 Icon text NOT NULL COMMENT '',
 Title varchar(100) unique NOT NULL  DEFAULT '' COMMENT '',
 ApiProvider varchar(1023) NOT NULL DEFAULT '' COMMENT '',
 ApiOscoreUrlKey varchar(100) unique NOT NULL COMMENT '',
 ApiUrl varchar(1023) NOT NULL  DEFAULT '' COMMENT '',
 Price varchar(100) NOT NULL  DEFAULT '' COMMENT '',
 ApiDesc varchar(1023) NOT NULL  DEFAULT '' COMMENT '',
 ErrorDesc TEXT(8191) NOT NULL COMMENT '',
 Specifications INT NOT NULL  DEFAULT 0 COMMENT '规格',
 Popularity INT NOT NULL DEFAULT 0 COMMENT '流行度',
 Delay INT NOT NULL DEFAULT 0 COMMENT '',
 SuccessRate INT NOT NULL DEFAULT 0 COMMENT '',
 InvokeFrequency BIGINT NOT NULL DEFAULT 0 COMMENT '',
 ApiState INT NOT NULL DEFAULT 0 COMMENT '',
 RequestType varchar(20) NOT NULL COMMENT '',
 Mark varchar(100) NOT NULL DEFAULT '' COMMENT '',
 ResponseParam TEXT(8191) NOT NULL  COMMENT '',
 ResponseExample TEXT(8191) NOT NULL COMMENT '',
 DataDesc varchar(255) NOT NULL DEFAULT '' COMMENT '',
 DataSource varchar(255) NOT NULL DEFAULT ''  COMMENT '',
 ApplicationScenario varchar(255) NOT NULL DEFAULT '' COMMENT '',
 ApiKind INT NOT NULL DEFAULT 1 COMMENT '',
 OntId varchar(50) NOT NULL DEFAULT '' COMMENT '',
 UserId varchar(255) NOT NULL DEFAULT '' COMMENT '',
 Author varchar(50) NOT NULL DEFAULT '' COMMENT '',
 ResourceId varchar(255) NOT NULL DEFAULT '',
 TokenHash char(255) NOT NULL DEFAULT '',
 OwnerAddress varchar(255) NOT NULL DEFAULT '',
 Abstract varchar(1023) NOT NULL DEFAULT '',
 NotifyDelete INT NOT NULL DEFAULT 0 COMMENT '',
 CreateTime BIGINT NOT NULL DEFAULT 0,
 UpdateTime BIGINT NOT NULL DEFAULT 0,
 PRIMARY KEY (ApiId),
 INDEX(Price),
 INDEX(Title),
 INDEX(ApiDesc),
 INDEX(ApiState),
 INDEX(OntId),
 INDEX(Author),
 INDEX(ApiKind)
)DEFAULT charset=utf8;

create table tbl_specifications
(
 Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 ApiId INT NOT NULL,
 Price  varchar(50) NOT NULL DEFAULT '' COMMENT '',
 Amount BIGINT NOT NULL DEFAULT 0,
 EffectiveDuration INT NOT NULL DEFAULT 0,
 SpecType INT NOT NULL DEFAULT 0,
 PRIMARY KEY (Id),
 CONSTRAINT FK_specifications_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId)
)DEFAULT charset=utf8;

create table tbl_request_param (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiId INT NOT NULL,
  ParamName varchar(50) NOT NULL DEFAULT '',
  Required  TINYINT NOT NULL,
  ParamWhere INT NOT NULL DEFAULT 0,
  ParamType varchar(10) NOT NULL DEFAULT '',
  Note varchar(255) NOT NULL DEFAULT '',
  ValueDesc varchar(50) NOT NULL DEFAULT '',
  ParamTag INT NOT NULL DEFAULT 0,
  PRIMARY KEY (Id),
  CONSTRAINT FK_request_param_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId),
  INDEX(ApiId)
)DEFAULT charset=utf8;

create table tbl_error_code (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ErrorCode INT NOT NULL,
  ErrorDesc varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (Id)
)DEFAULT charset=utf8;

create table tbl_order (
  OrderId varchar(100) unique NOT NULL COMMENT '',
  Title varchar(100) NOT NULL COMMENT '',
  ProductName varchar(50) NOT NULL DEFAULT '' COMMENT '',
  OrderType varchar(50) NOT NULL DEFAULT ''  COMMENT '',
  OrderTime BIGINT NOT NULL DEFAULT 0 COMMENT '下单时间',
  PayTime  BIGINT NOT NULL DEFAULT 0  COMMENT '支付时间',
  State INT NOT NULL DEFAULT 0,
  NotifyPay INT NOT NULL DEFAULT 0,
  Amount varchar(255) NOT NULL DEFAULT '' COMMENT '',
  OntId varchar(50) NOT NULL DEFAULT '' COMMENT '用户ontid',
  UserId varchar(255) NOT NULL DEFAULT '' COMMENT '',
  UserName varchar(50) NOT NULL DEFAULT '' COMMENT '',
  TxHash varchar(255) NOT NULL DEFAULT '' COMMENT '',
  Price varchar(50) NOT NULL DEFAULT ''  COMMENT '',
  ApiId INT NOT NULL COMMENT '',
  ToolBoxId INT NOT NULL COMMENT '',
  ApiUrl varchar(255) NOT NULL  DEFAULT '' COMMENT '',
  SpecificationsId INT NOT NULL COMMENT '规格',
  OrderKind INT NOT NULL COMMENT '',
  Request varchar(4095) NOT NULL COMMENT '币种',
  Result varchar(4095) NOT NULL COMMENT '币种',
  ApiKey varchar(50) NOT NULL DEFAULT '',
  Reason varchar(1023) NOT NULL DEFAULT '',
  PRIMARY KEY (OrderId),
  CONSTRAINT FK_tbl_order_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId),
  INDEX(OntId)
)DEFAULT charset=utf8;


create table tbl_api_key (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiKey varchar(50) unique NOT NULL  DEFAULT '',
  ApiId INT NOT NULL,
  OrderId varchar(100) unique NOT NULL COMMENT '',
  RequestLimit BIGINT NOT NULL DEFAULT 0,
  UsedNum BIGINT NOT NULL DEFAULT 0,
  OntId varchar(50) NOT NULL DEFAULT '',
  UserId varchar(255) NOT NULL DEFAULT '',
  OutDate BIGINT NOT NULL DEFAULT 0,
  ApiKeyType INT NOT NULL DEFAULT 0,
  CreateTime BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (Id),
  foreign key(OrderId) references tbl_order(OrderId),
  foreign key(ApiId) references tbl_api_basic_info(ApiId),
  INDEX(ApiKey),
  INDEX(OntId)
)DEFAULT charset=utf8;

create table tbl_api_test_key (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiKey varchar(50) unique NOT NULL  DEFAULT '',
  ApiId INT NOT NULL,
  OrderId varchar(20) NOT NULL DEFAULT 'TST_ORDER' COMMENT '',
  RequestLimit BIGINT NOT NULL DEFAULT 0,
  UsedNum BIGINT NOT NULL DEFAULT 0,
  OntId varchar(50) NOT NULL DEFAULT '',
  OutDate BIGINT NOT NULL DEFAULT 0,
  ApiKeyType INT NOT NULL DEFAULT 0,
  CreateTime BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (Id),
  foreign key(ApiId) references tbl_api_basic_info(ApiId),
  INDEX(ApiId),
  INDEX(ApiKey),
  INDEX(OntId)
) DEFAULT charset=utf8;

CREATE TABLE `tbl_qr_code` (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  QrCodeId varchar(100) unique NOT NULL DEFAULT '',
  Ver varchar(50) NOT NULL DEFAULT '',
  OrderId varchar(100) NOT NULL DEFAULT '' ,
  Requester varchar(50) NOT NULL DEFAULT '',
  Signature varchar(200) NOT NULL DEFAULT '',
  Signer varchar(50) NOT NULL DEFAULT '',
  QrCodeData text,
  Callback varchar(400) NOT NULL DEFAULT '',
  Exp BIGINT NOT NULL DEFAULT 0,
  Chain varchar(50) NOT NULL DEFAULT '',
  QrCodeDesc varchar(100) NOT NULL DEFAULT '',
  ContractType varchar(10) NOT NULL DEFAULT '',
  PRIMARY KEY (Id),
  foreign key(OrderId) references tbl_order(OrderId),
  INDEX(QrCodeId)
)DEFAULT charset=utf8;

CREATE TABLE IF NOT EXISTS `tbl_country_city` (
	Id INT NOT NULL AUTO_INCREMENT,
	Country varchar(50) NOT NULL,
	City varchar(50) UNIQUE NOT NULL,
	Lat varchar(50) NOT NULL,
	Lng varchar(50) NOT NULL,
	PRIMARY KEY (Id)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_algorithm` (
	Id INT NOT NULL AUTO_INCREMENT,
	AlgName varchar(255) UNIQUE NOT NULL,
	Provider varchar(255) NOT NULL DEFAULT '',
	Description varchar(255) NOT NULL DEFAULT '',
	Price varchar(255) NOT NULL DEFAULT '' COMMENT '',
	Coin varchar(20) NOT NULL COMMENT '币种',
	ResourceId varchar(255) UNIQUE NOT NULL,
	TokenHash char(255) UNIQUE NOT NULL,
	OwnerAddress varchar(255) UNIQUE NOT NULL,
	State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
	CreateTime TIMESTAMP DEFAULT current_timestamp,
	PRIMARY KEY(Id)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_env` (
	Id INT NOT NULL AUTO_INCREMENT,
	EnvName varchar(255) UNIQUE NOT NULL,
	Provider varchar(255) NOT NULL DEFAULT '',
	Description varchar(255) NOT NULL DEFAULT '',
	Price varchar(255) NOT NULL DEFAULT '0' COMMENT '',
	Coin varchar(20) NOT NULL COMMENT '币种',
	ServiceUrl varchar(255) NOT NULL,
	ResourceId varchar(255) UNIQUE NOT NULL,
	TokenHash char(255) UNIQUE NOT NULL,
	OwnerAddress varchar(255) UNIQUE NOT NULL,
	State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
	CreateTime TIMESTAMP DEFAULT current_timestamp,
	PRIMARY KEY(Id)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_api_algorithm` (
	Id INT NOT NULL AUTO_INCREMENT,
	ApiId INT NOT NULL,
	AlgorithmId INT NOT NULL,
	State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
	CreateTime TIMESTAMP DEFAULT current_timestamp,
	foreign key(ApiId) references tbl_api_basic_info(ApiId),
	foreign key(AlgorithmId) references tbl_algorithm(Id),
	PRIMARY KEY(Id)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_algorithm_env` (
	Id INT NOT NULL AUTO_INCREMENT,
	AlgorithmId INT NOT NULL,
	EnvId INT NOT NULL,
	State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
	CreateTime TIMESTAMP DEFAULT current_timestamp,
	foreign key(EnvId) references tbl_env(Id),
	foreign key(AlgorithmId) references tbl_algorithm(Id),
	PRIMARY KEY(Id)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_tool_box` (
	Id INT NOT NULL AUTO_INCREMENT,
	Title varchar(255) UNIQUE NOT NULL COMMENT 'comresponse to api_basic_info.ApiType',
	ToolBoxDesc  varchar(255) NOT NULL,
	ToolBoxType varchar(255) NOT NULL,
	Coin varchar(255) NOT NULL,
	Price varchar(255) NOT NULL,
	Icon varchar(255),
	State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
	CreateTime TIMESTAMP DEFAULT current_timestamp,
	PRIMARY KEY(Id)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_order_api_comment` (
	Id INT NOT NULL AUTO_INCREMENT,
	OrderId varchar(255) unique NOT NULL COMMENT '',
	ApiId INT NOT NULL COMMENT '',
	ToolBoxId INT NOT NULL COMMENT '',
	StarNum INT NOT NULL COMMENT '',
	Comments varchar(255),
	UserName varchar(255) NOT NULL DEFAULT '',
	OntId varchar(50) NOT NULL DEFAULT '' COMMENT '',
	CommentTime BIGINT NOT NULL DEFAULT 0 COMMENT '',
	State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
	INDEX(ApiId),
	INDEX(ToolBoxId),
	PRIMARY KEY (Id),
	foreign key(OrderId) references tbl_order(OrderId),
	foreign key(ApiId) references tbl_api_basic_info(ApiId)
)DEFAULT charset=utf8;

CREATE TABLE `tbl_score` (
	OntId varchar(50) unique NOT NULL DEFAULT '' COMMENT '',
	TotalCommentNum BIGINT NOT NULL DEFAULT 0,
	TotalScore BIGINT NOT NULL DEFAULT 0
)DEFAULT charset=utf8;

CREATE TABLE `tbl_api_crypto_info` (
	ApiId INT unique NOT NULL COMMENT '',
	CryptoType INT NOT NULL DEFAULT 0 COMMENT '',
	CryptoWhere INT NOT NULL DEFAULT 0 COMMENT '',
	CryptoKey varchar(255) NOT NULL DEFAULT '',
	CryptoValue varchar(255) NOT NULL DEFAULT '',
	foreign key(ApiId) references tbl_api_basic_info(ApiId)
)DEFAULT charset=utf8;

