/*
 * Named proxy callback function, see 
 *
 * https://code.google.com/p/go-wiki/wiki/cgo
 */

#include "callback_proxy.h"

int go_readproxy(void *p, unsigned char *buf, int nbytes);

int
c_readproxy(void *p, unsigned char *buf, int nbytes)
{
	return go_readproxy(p, buf, nbytes);
}
