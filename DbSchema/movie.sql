-- phpMyAdmin SQL Dump
-- version 4.7.1
-- https://www.phpmyadmin.net/
--
-- 主機: 127.0.0.1
-- 產生時間： 2017 年 06 月 04 日 11:34
-- 伺服器版本: 5.7.17
-- PHP 版本： 5.6.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 資料庫： `movie`
--

-- --------------------------------------------------------

--
-- 資料表結構 `movieList`
--

CREATE TABLE `movieList` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `cname` varchar(128) DEFAULT NULL,
  `ename` varchar(128) DEFAULT NULL,
  `releaseTime` timestamp NULL DEFAULT NULL,
  `type` varchar(128) DEFAULT NULL,
  `duration` varchar(128) DEFAULT NULL,
  `director` varchar(128) DEFAULT NULL,
  `actor` text,
  `company` varchar(128) DEFAULT NULL,
  `website` varchar(256) DEFAULT NULL,
  `score` double DEFAULT NULL,
  `intro` text,
  `imgPath` varchar(64) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 已匯出資料表的索引
--

--
-- 資料表索引 `movieList`
--
ALTER TABLE `movieList`
  ADD PRIMARY KEY (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
