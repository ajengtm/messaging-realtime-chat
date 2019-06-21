CREATE DATABASE messaging;

CREATE TABLE `messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sender_id` int(11) NOT NULL,
  `recepient_id` int(11) NOT NULL,
  `message` text NOT NULL,
  `status` tinyint(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `messages` (`id`, `sender_id`, `recepient_id`, `message`, `status`, `created_at`, `updated_at`)
VALUES
	(82,1,2,'yayayayaya',4,'2019-05-20 06:21:30','2019-05-20 06:21:39'),
	(83,1,2,'yayayayaya',4,'2019-05-20 06:21:31','2019-05-20 06:21:42'),
	(84,1,2,'yayayayaya',4,'2019-05-20 06:21:33','2019-05-20 06:21:45'),
	(85,1,2,'yayayayaya',3,'2019-05-20 12:46:19','2019-05-20 12:48:20');