在上一节中介绍了 socket 的 `Listen` 方法，这里进一步介绍 `Accept` 和 `Read`，`Write` 方法。

# 1. Accept
Accept 的核心逻辑在于：
```
func (ln *TCPListener) accept() (*TCPConn, error) {
	fd, err := ln.fd.accept()
	if err != nil {
		return nil, err
	}
	tc := newTCPConn(fd)
	if ln.lc.KeepAlive >= 0 {
		setKeepAlive(fd, true)
		ka := ln.lc.KeepAlive
		if ln.lc.KeepAlive == 0 {
			ka = defaultTCPKeepAlive
		}
		setKeepAlivePeriod(fd, ka)
	}
	return tc, nil
}
```

通过 socket 返回的 fd 调用 `accept` 方法从 socket 上接收数据。accept 返回新 fd，通过该新 fd 建立 tcp 连接。并且通过 `setKeepAlive` 和 `setKeepAlivePeriod` 函数添加对应该新 fd 的 KeepAlive 属性：`tcp_keepalive_time`, `tcp_keepalive_intvl` 和 `tcp_keepalive_probes`。

在 KeepAlive 函数中有一段函数 `runtime.KeepAlive` 比较有意思：
```
func setKeepAlive(fd *netFD, keepalive bool) error {
	err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, boolint(keepalive))
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}
```

它的存在是为了让 fd 不会被 GC 回收，更多信息可参考 [issue_21402](https://github.com/golang/go/issues/2140) 和 [go 变量逃逸分析](https://www.cnblogs.com/xingzheanan/p/16082035.html)。

继续看 `accept` 方法：
```
func (fd *netFD) accept() (netfd *netFD, err error) {
	d, rsa, errcall, err := fd.pfd.Accept()
	...
	if netfd, err = newFD(d, fd.family, fd.sotype, fd.net); err != nil {
		poll.CloseFunc(d)
		return nil, err
	}

	netfd.setAddr(netfd.addrFunc()(lsa), netfd.addrFunc()(rsa))
	return netfd, nil
}
```

`netFD` 包的是 pfd `poll.FD`，调用 pfd 的 `Accept` 方法返回 socket 上的系统文件描述符 `d`。将 `d` 包装成 `netfd`，接着通过 `setAddr` 设置 netfd 的本地地址 laddr 和 client 端地址 raddr。

`poll.FD` 的 `Accept` 是重头戏了，接着看：
```
func (fd *FD) Accept() (int, syscall.Sockaddr, string, error) {
	for {
		s, rsa, errcall, err := accept(fd.Sysfd)
		if err == nil {
			return s, rsa, "", err
		}
		switch err {
		case syscall.EINTR:
			continue
		case syscall.EAGAIN:
			if fd.pd.pollable() {
				if err = fd.pd.waitRead(fd.isFile); err == nil {
					continue
				}
			}
		...
	}
}

func accept(s int) (int, syscall.Sockaddr, string, error) {
	ns, sa, err := Accept4Func(s, syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC)
	switch err {
	case nil:
		return ns, sa, "", nil
	...
    }
}

func accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error) {
	r0, _, e1 := Syscall6(SYS_ACCEPT4, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), uintptr(flags), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
```

poll.FD 的 `accept` 方法中做了下面几件事：
1. `accept` 函数经过 `Accept4Func`, `accept4` 到系统调用，通过系统调用号 `SYS_ACCEPT4` 和文件描述符 `fd.Sysfd` 返回作用在 socket 上的系统文件描述符和远端 socket 地址。
2. 这里 accept 是非阻塞的，意味着即使没有 client 连接也会返回。此时返回的 err 类型为 `syscall.EAGAIN`。
3. 进入到 EAGAIN 错误类型中，会通过 `fd.pd.pollable` 方法判断是否为 `true`。如果为 `true` 阻塞当前 goroutine 直到有新的可读数据。

Accept 的实现简单介绍基本告一段落了，下面继续看 socket 的 `Read` 和 `Write` 实现。

# 2. Read 和 Write
## 2.1 Read
`Read` 经过层层调用到 poll.FD 的 Read 方法：
```
func (fd *FD) Read(p []byte) (int, error) {
	...
	if err := fd.pd.prepareRead(fd.isFile); err != nil {
		return 0, err
	}
	if fd.IsStream && len(p) > maxRW {
		p = p[:maxRW]
	}
	for {
		n, err := ignoringEINTRIO(syscall.Read, fd.Sysfd, p)
		if err != nil {
			n = 0
			if err == syscall.EAGAIN && fd.pd.pollable() {
				if err = fd.pd.waitRead(fd.isFile); err == nil {
					continue
				}
			}
		}
		err = fd.eofError(n, err)
		return n, err
	}
}
```

从上述代码可以发现：
- 网络处理逻辑通过层层封装走到 poll 的 `Read`，poll 是不区分文件还是网络数据的。因此，在 `prepareRead` 中需要通过 `fd.isFile` 判断。
- `maxRW` 是能读取数据的最大字节，这里是 1G。原因分析在注释中：
```
// Darwin and FreeBSD can't read or write 2GB+ files at a time,
// even on 64-bit systems.
// The same is true of socket implementations on many systems.
// See golang.org/issue/7812 and golang.org/issue/16266.
// Use 1GB instead of, say, 2GB-1, to keep subsequent reads aligned.
```
- `ignoringEINTRIO` 中通过 `syscall.Read` 函数，作用在系统调用上，通过系统调用号和文件描述符 `fd.Sysfd` 读取 socket 的数据到 p。
- 类似 Accept，如果 `ignoringEINTRIO` 返回错误 `syscall.EAGAIN`，并且 `fd.pd.pollable` 是 true 的话，会阻塞当前 goroutine 等待读取数据。
- 进入到 `eofError` 逻辑。对于文件，如果读到 `EOF` 则说明文件结束。对于网络数据，err 返回为 nil。

## Write
类似于 `Read`，`Write` 的核心逻辑在：
```
// Write implements io.Writer.
func (fd *FD) Write(p []byte) (int, error) {
	...
	for {
		max := len(p)
		if fd.IsStream && max-nn > maxRW {
			max = nn + maxRW
		}
		n, err := ignoringEINTRIO(syscall.Write, fd.Sysfd, p[nn:max])
		if n > 0 {
			nn += n
		}
		if nn == len(p) {
			return nn, err
		}
		if err == syscall.EAGAIN && fd.pd.pollable() {
			if err = fd.pd.waitWrite(fd.isFile); err == nil {
				continue
			}
		}
		if err != nil {
			return nn, err
		}
		if n == 0 {
			return nn, io.ErrUnexpectedEOF
		}
	}
}
```

通过 `syscall.Write` 函数进入系统调用，执行 `Write` 调用作用于系统文件描述符 `fd.Sysfd` 写数据到 `p`。如果返回 `EAGAIN` 且 `pollable` 为 `true` 的话则阻塞当前 goroutine 进入 `waitWrite`。直到数据写完，跳出 for 循环。
