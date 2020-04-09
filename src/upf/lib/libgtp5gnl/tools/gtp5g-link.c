/* Command line utility to create GTP link */

/* (C) 2014 by sysmocom - s.f.m.c. GmbH
 * (C) 2016 by Pablo Neira Ayuso <pablo@netfilter.org>
 * (C) 2017 by Pablo Neira Ayuso <pablo@gnumonks.org>
 *
 * Author: Yao-Wen Chang <yaowenowo@gmail.com>
 *
 * All Rights Reserved
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <time.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include <libmnl/libmnl.h>
#include <linux/if.h>
#include <linux/if_link.h>
#include <linux/rtnetlink.h>

#include <linux/gtp5g.h>
#include <linux/if_link.h>

#include <libgtp5gnl/gtp5gnl.h>

int main(int argc, char *argv[])
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	int ret, ran_mode = 0;

	if (argc < 3) {
		printf("Usage: %s <add|del> <device>\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	if (!strcmp(argv[1], "del")) {
		printf("destroying gtp interface...\n");
		if (gtp_dev_destroy(argv[2]) < 0)
			perror("gtp_dev_destroy");

		return 0;
	}

	if (argc > 3 && !strcmp(argv[3], "--ran"))
		ran_mode = 1;

	int fd = socket(AF_INET, SOCK_DGRAM, 0);
	struct sockaddr_in sockaddr_fd = {
		.sin_family	= AF_INET,
		.sin_port	= htons(2152),
		.sin_addr	= {
			.s_addr 	= INADDR_ANY,
		},
	};

	if (bind(fd, (struct sockaddr *) &sockaddr_fd,
		 sizeof(sockaddr_fd)) < 0) {
		perror("bind");
		exit(EXIT_FAILURE);
	}

	if (ran_mode)
		ret = gtp_dev_create_ran(-1, argv[2], fd);
	else
		ret = gtp_dev_create(-1, argv[2], fd);
	if (ret < 0) {
		perror("cannot create 5G GTP device\n");
		exit(EXIT_FAILURE);
	}

	fprintf(stderr, "WARNING: attaching dummy socket descriptors. Keep "
			"this process running for testing purposes.\n");

	while (1) {
		struct sockaddr_in addr;
		socklen_t len = sizeof(addr);

		ret = recvfrom(fd, buf, sizeof(buf), 0,
			       (struct sockaddr *)&addr, &len);
		printf("received %d bytes via UDP socket\n", ret);
	}

	return 0;
}
