create table category
(
  id character varying(100) not null,
  data json[],
  primary key (id)
);

create table channel
(
  id character varying(100) not null,
  count integer,
  country character varying(25),
  customurl character varying(250),
  description character varying,
  favorites character varying,
  highthumbnail character varying(255),
  itemcount integer,
  likes character varying,
  localizeddescription character varying,
  localizedtitle character varying(255),
  mediumthumbnail character varying(255),
  playlistcount integer,
  playlistitemcount integer,
  playlistvideocount integer,
  playlistvideoitemcount integer,
  publishedat timestamp with time zone,
  thumbnail character varying(255),
  lastupload timestamp with time zone,
  title character varying(255),
  uploads character varying(100),
  channels character varying[],
  primary key (id)
);

create table channelsync
(
  id character varying(100) not null,
  synctime timestamp with time zone,
  uploads character varying(100),
  primary key (id)
);

create table playlist
(
  id character varying(100) not null,
  channelid character varying(100),
  channeltitle character varying(255),
  count integer,
  itemcount integer,
  description character varying,
  highthumbnail character varying(255),
  localizeddescription character varying,
  localizedtitle character varying(255),
  maxresthumbnail character varying(255),
  mediumthumbnail character varying(255),
  publishedat timestamp with time zone,
  standardthumbnail character varying(255),
  thumbnail character varying(255),
  title character varying(255),
  primary key (id)
);

create table playlistvideo
(
  id character varying(100) not null,
  videos character varying[],
  primary key (id)
);

create table video
(
  id character varying(100) not null,
  caption character varying(255),
  categoryid character varying(20),
  channelid character varying(100),
  channeltitle character varying(255),
  defaultaudiolanguage character varying(255),
  defaultlanguage character varying(255),
  definition smallint,
  description character varying,
  dimension character varying(20),
  duration bigint,
  highthumbnail character varying(255),
  licensedcontent boolean,
  livebroadcastcontent character varying(255),
  localizeddescription character varying,
  localizedtitle character varying(255),
  maxresthumbnail character varying(255),
  mediumthumbnail character varying(255),
  projection character varying(255),
  publishedat timestamp with time zone,  standardthumbnail character varying(255),
  tags character varying[],
  thumbnail character varying(255),
  title character varying(255),
  blockedregions character varying(100)[],
  allowedregions character varying(100)[],
  primary key (id)
)
