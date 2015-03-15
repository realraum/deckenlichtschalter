##
##  spreadspace avr projects
##
##
##  Copyright (C) 2013-2015 Christian Pointner <equinox@spreadspace.org>
##
##  This file is part of spreadspace avr projects.
##
##  spreadspace avr projects is free software: you can redistribute it and/or modify
##  it under the terms of the GNU General Public License as published by
##  the Free Software Foundation, either version 3 of the License, or
##  any later version.
##
##  spreadspace avr projects is distributed in the hope that it will be useful,
##  but WITHOUT ANY WARRANTY; without even the implied warranty of
##  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
##  GNU General Public License for more details.
##
##  You should have received a copy of the GNU General Public License
##  along with spreadspace avr projects. If not, see <http://www.gnu.org/licenses/>.
##

all: update-spreadavr

.spreadavr.prepared:
	git clone git://git.spreadspace.org/avr-utils.git
	touch $@

update-spreadavr: .spreadavr.prepared
	cd avr-utils/; git pull
	$(MAKE) -C avr-utils/contrib/
