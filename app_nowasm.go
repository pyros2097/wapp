// +build !wasm

package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/markbates/pkger"
)

func Run() {
	isLambda := os.Getenv("_LAMBDA_SERVER_PORT") != ""
	if !isLambda {
		wd, err := os.Getwd()
		if err != nil {
			println("could not get wd")
			return
		}
		assetsFS := http.FileServer(pkger.Dir(filepath.Join(wd, "assets")))
		AppRouter.GET("/assets/*filepath", func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.Replace(r.URL.Path, "/assets", "", 1)
			assetsFS.ServeHTTP(w, r)
		})
	}
	if isLambda {
		println("running in lambda mode")
		lambda.Start(AppRouter.Lambda)
	} else {
		println("Serving on HTTP port: 1234")
		http.ListenAndServe(":1234", AppRouter)
	}
}

func Reload() {
	panic("wasm required")
}

func Route(path string, render RenderFunc) {
	println("registering route: " + path)
	AppRouter.GET(path, render)
}

func writePage(ui UI, w io.Writer) {
	// "<!DOCTYPE html>\n"
	Html(
		Head(
			Title(helmet.Title),
			Meta("description", helmet.Description),
			Meta("author", helmet.Author),
			Meta("keywords", helmet.Keywords),
			Meta("viewport", "width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0, viewport-fit=cover"),
			Link("icon", "/assets/icon.png"),
			Link("apple-touch-icon", "/assets/icon.png"),
			Link("stylesheet", "/assets/styles.css"),
			Script(`
(() => {
	// Map multiple JavaScript environments to a single common API,
	// preferring web standards over Node.js API.
	//
	// Environments considered:
	// - Browsers
	// - Node.js
	// - Electron
	// - Parcel

	if (typeof global !== "undefined") {
		// global already exists
	} else if (typeof window !== "undefined") {
		window.global = window;
	} else if (typeof self !== "undefined") {
		self.global = self;
	} else {
		throw new Error("cannot export Go (neither global, window nor self is defined)");
	}

	if (!global.require && typeof require !== "undefined") {
		global.require = require;
	}

	if (!global.fs && global.require) {
		global.fs = require("fs");
	}

	const enosys = () => {
		const err = new Error("not implemented");
		err.code = "ENOSYS";
		return err;
	};

	if (!global.fs) {
		let outputBuf = "";
		global.fs = {
			constants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused
			writeSync(fd, buf) {
				outputBuf += decoder.decode(buf);
				const nl = outputBuf.lastIndexOf("\n");
				if (nl != -1) {
					console.log(outputBuf.substr(0, nl));
					outputBuf = outputBuf.substr(nl + 1);
				}
				return buf.length;
			},
			write(fd, buf, offset, length, position, callback) {
				if (offset !== 0 || length !== buf.length || position !== null) {
					callback(enosys());
					return;
				}
				const n = this.writeSync(fd, buf);
				callback(null, n);
			},
			chmod(path, mode, callback) { callback(enosys()); },
			chown(path, uid, gid, callback) { callback(enosys()); },
			close(fd, callback) { callback(enosys()); },
			fchmod(fd, mode, callback) { callback(enosys()); },
			fchown(fd, uid, gid, callback) { callback(enosys()); },
			fstat(fd, callback) { callback(enosys()); },
			fsync(fd, callback) { callback(null); },
			ftruncate(fd, length, callback) { callback(enosys()); },
			lchown(path, uid, gid, callback) { callback(enosys()); },
			link(path, link, callback) { callback(enosys()); },
			lstat(path, callback) { callback(enosys()); },
			mkdir(path, perm, callback) { callback(enosys()); },
			open(path, flags, mode, callback) { callback(enosys()); },
			read(fd, buffer, offset, length, position, callback) { callback(enosys()); },
			readdir(path, callback) { callback(enosys()); },
			readlink(path, callback) { callback(enosys()); },
			rename(from, to, callback) { callback(enosys()); },
			rmdir(path, callback) { callback(enosys()); },
			stat(path, callback) { callback(enosys()); },
			symlink(path, link, callback) { callback(enosys()); },
			truncate(path, length, callback) { callback(enosys()); },
			unlink(path, callback) { callback(enosys()); },
			utimes(path, atime, mtime, callback) { callback(enosys()); },
		};
	}

	if (!global.process) {
		global.process = {
			getuid() { return -1; },
			getgid() { return -1; },
			geteuid() { return -1; },
			getegid() { return -1; },
			getgroups() { throw enosys(); },
			pid: -1,
			ppid: -1,
			umask() { throw enosys(); },
			cwd() { throw enosys(); },
			chdir() { throw enosys(); },
		}
	}

	if (!global.crypto) {
		const nodeCrypto = require("crypto");
		global.crypto = {
			getRandomValues(b) {
				nodeCrypto.randomFillSync(b);
			},
		};
	}

	if (!global.performance) {
		global.performance = {
			now() {
				const [sec, nsec] = process.hrtime();
				return sec * 1000 + nsec / 1000000;
			},
		};
	}

	if (!global.TextEncoder) {
		global.TextEncoder = require("util").TextEncoder;
	}

	if (!global.TextDecoder) {
		global.TextDecoder = require("util").TextDecoder;
	}

	// End of polyfills for common API.

	const encoder = new TextEncoder("utf-8");
	const decoder = new TextDecoder("utf-8");
	var logLine = [];

	global.Go = class {
		constructor() {
			this._callbackTimeouts = new Map();
			this._nextCallbackTimeoutID = 1;

			const mem = () => {
				// The buffer may change when requesting more memory.
				return new DataView(this._inst.exports.memory.buffer);
			}

			const setInt64 = (addr, v) => {
				mem().setUint32(addr + 0, v, true);
				mem().setUint32(addr + 4, Math.floor(v / 4294967296), true);
			}

			const getInt64 = (addr) => {
				const low = mem().getUint32(addr + 0, true);
				const high = mem().getInt32(addr + 4, true);
				return low + high * 4294967296;
			}

			const loadValue = (addr) => {
				const f = mem().getFloat64(addr, true);
				if (f === 0) {
					return undefined;
				}
				if (!isNaN(f)) {
					return f;
				}

				const id = mem().getUint32(addr, true);
				return this._values[id];
			}

			const storeValue = (addr, v) => {
				const nanHead = 0x7FF80000;

				if (typeof v === "number") {
					if (isNaN(v)) {
						mem().setUint32(addr + 4, nanHead, true);
						mem().setUint32(addr, 0, true);
						return;
					}
					if (v === 0) {
						mem().setUint32(addr + 4, nanHead, true);
						mem().setUint32(addr, 1, true);
						return;
					}
					mem().setFloat64(addr, v, true);
					return;
				}

				switch (v) {
					case undefined:
						mem().setFloat64(addr, 0, true);
						return;
					case null:
						mem().setUint32(addr + 4, nanHead, true);
						mem().setUint32(addr, 2, true);
						return;
					case true:
						mem().setUint32(addr + 4, nanHead, true);
						mem().setUint32(addr, 3, true);
						return;
					case false:
						mem().setUint32(addr + 4, nanHead, true);
						mem().setUint32(addr, 4, true);
						return;
				}

				let id = this._ids.get(v);
				if (id === undefined) {
					id = this._idPool.pop();
					if (id === undefined) {
						id = this._values.length;
					}
					this._values[id] = v;
					this._goRefCounts[id] = 0;
					this._ids.set(v, id);
				}
				this._goRefCounts[id]++;
				let typeFlag = 1;
				switch (typeof v) {
					case "string":
						typeFlag = 2;
						break;
					case "symbol":
						typeFlag = 3;
						break;
					case "function":
						typeFlag = 4;
						break;
				}
				mem().setUint32(addr + 4, nanHead | typeFlag, true);
				mem().setUint32(addr, id, true);
			}

			const loadSlice = (array, len, cap) => {
				return new Uint8Array(this._inst.exports.memory.buffer, array, len);
			}

			const loadSliceOfValues = (array, len, cap) => {
				const a = new Array(len);
				for (let i = 0; i < len; i++) {
					a[i] = loadValue(array + i * 8);
				}
				return a;
			}

			const loadString = (ptr, len) => {
				return decoder.decode(new DataView(this._inst.exports.memory.buffer, ptr, len));
			}

			const timeOrigin = Date.now() - performance.now();
			this.importObject = {
				wasi_unstable: {
					// https://github.com/bytecodealliance/wasmtime/blob/master/docs/WASI-api.md#__wasi_fd_write
					fd_write: function(fd, iovs_ptr, iovs_len, nwritten_ptr) {
						let nwritten = 0;
						if (fd == 1) {
							for (let iovs_i=0; iovs_i<iovs_len;iovs_i++) {
								let iov_ptr = iovs_ptr+iovs_i*8; // assuming wasm32
								let ptr = mem().getUint32(iov_ptr + 0, true);
								let len = mem().getUint32(iov_ptr + 4, true);
								for (let i=0; i<len; i++) {
									let c = mem().getUint8(ptr+i);
									if (c == 13) { // CR
										// ignore
									} else if (c == 10) { // LF
										// write line
										let line = decoder.decode(new Uint8Array(logLine));
										logLine = [];
										console.log(line);
									} else {
										logLine.push(c);
									}
								}
							}
						} else {
							console.error('invalid file descriptor:', fd);
						}
						mem().setUint32(nwritten_ptr, nwritten, true);
						return 0;
					},
				},
				env: {
					// func ticks() float64
					"runtime.ticks": () => {
						return timeOrigin + performance.now();
					},

					// func sleepTicks(timeout float64)
					"runtime.sleepTicks": (timeout) => {
						// Do not sleep, only reactivate scheduler after the given timeout.
						setTimeout(this._inst.exports.go_scheduler, timeout);
					},

					// func finalizeRef(v ref)
					"syscall/js.finalizeRef": (sp) => {
						// Note: TinyGo does not support finalizers so this should never be
						// called.
						console.error('syscall/js.finalizeRef not implemented');
					},

					// func stringVal(value string) ref
					"syscall/js.stringVal": (ret_ptr, value_ptr, value_len) => {
						const s = loadString(value_ptr, value_len);
						storeValue(ret_ptr, s);
					},

					// func valueGet(v ref, p string) ref
					"syscall/js.valueGet": (retval, v_addr, p_ptr, p_len) => {
						let prop = loadString(p_ptr, p_len);
						let value = loadValue(v_addr);
						let result = Reflect.get(value, prop);
						storeValue(retval, result);
					},

					// func valueSet(v ref, p string, x ref)
					"syscall/js.valueSet": (v_addr, p_ptr, p_len, x_addr) => {
						const v = loadValue(v_addr);
						const p = loadString(p_ptr, p_len);
						const x = loadValue(x_addr);
						Reflect.set(v, p, x);
					},

					// func valueDelete(v ref, p string)
					"syscall/js.valueDelete": (v_addr, p_ptr, p_len) => {
						const v = loadValue(v_addr);
						const p = loadString(p_ptr, p_len);
						Reflect.deleteProperty(v, p);
					},

					// func valueIndex(v ref, i int) ref
					"syscall/js.valueIndex": (ret_addr, v_addr, i) => {
						storeValue(ret_addr, Reflect.get(loadValue(v_addr), i));
					},

					// valueSetIndex(v ref, i int, x ref)
					"syscall/js.valueSetIndex": (v_addr, i, x_addr) => {
						Reflect.set(loadValue(v_addr), i, loadValue(x_addr));
					},

					// func valueCall(v ref, m string, args []ref) (ref, bool)
					"syscall/js.valueCall": (ret_addr, v_addr, m_ptr, m_len, args_ptr, args_len, args_cap) => {
						const v = loadValue(v_addr);
						const name = loadString(m_ptr, m_len);
						const args = loadSliceOfValues(args_ptr, args_len, args_cap);
						try {
							const m = Reflect.get(v, name);
							storeValue(ret_addr, Reflect.apply(m, v, args));
							mem().setUint8(ret_addr + 8, 1);
						} catch (err) {
							storeValue(ret_addr, err);
							mem().setUint8(ret_addr + 8, 0);
						}
					},

					// func valueInvoke(v ref, args []ref) (ref, bool)
					"syscall/js.valueInvoke": (ret_addr, v_addr, args_ptr, args_len, args_cap) => {
						try {
							const v = loadValue(v_addr);
							const args = loadSliceOfValues(args_ptr, args_len, args_cap);
							storeValue(ret_addr, Reflect.apply(v, undefined, args));
							mem().setUint8(ret_addr + 8, 1);
						} catch (err) {
							storeValue(ret_addr, err);
							mem().setUint8(ret_addr + 8, 0);
						}
					},

					// func valueNew(v ref, args []ref) (ref, bool)
					"syscall/js.valueNew": (ret_addr, v_addr, args_ptr, args_len, args_cap) => {
						const v = loadValue(v_addr);
						const args = loadSliceOfValues(args_ptr, args_len, args_cap);
						try {
							storeValue(ret_addr, Reflect.construct(v, args));
							mem().setUint8(ret_addr + 8, 1);
						} catch (err) {
							storeValue(ret_addr, err);
							mem().setUint8(ret_addr+ 8, 0);
						}
					},

					// func valueLength(v ref) int
					"syscall/js.valueLength": (v_addr) => {
						return loadValue(v_addr).length;
					},

					// valuePrepareString(v ref) (ref, int)
					"syscall/js.valuePrepareString": (ret_addr, v_addr) => {
						const s = String(loadValue(v_addr));
						const str = encoder.encode(s);
						storeValue(ret_addr, str);
						setInt64(ret_addr + 8, str.length);
					},

					// valueLoadString(v ref, b []byte)
					"syscall/js.valueLoadString": (v_addr, slice_ptr, slice_len, slice_cap) => {
						const str = loadValue(v_addr);
						loadSlice(slice_ptr, slice_len, slice_cap).set(str);
					},

					// func valueInstanceOf(v ref, t ref) bool
					//"syscall/js.valueInstanceOf": (sp) => {
					//	mem().setUint8(sp + 24, loadValue(sp + 8) instanceof loadValue(sp + 16));
					//},

					// func copyBytesToGo(dst []byte, src ref) (int, bool)
					"syscall/js.copyBytesToGo": (ret_addr, dest_addr, dest_len, dest_cap, source_addr) => {
						let num_bytes_copied_addr = ret_addr;
						let returned_status_addr = ret_addr + 4; // Address of returned boolean status variable

						const dst = loadSlice(dest_addr, dest_len);
						const src = loadValue(source_addr);
						if (!(src instanceof Uint8Array)) {
							mem().setUint8(returned_status_addr, 0); // Return "not ok" status
							return;
						}
						const toCopy = src.subarray(0, dst.length);
						dst.set(toCopy);
						setInt64(num_bytes_copied_addr, toCopy.length);
						mem().setUint8(returned_status_addr, 1); // Return "ok" status
					},

					// copyBytesToJS(dst ref, src []byte) (int, bool)
					// Originally copied from upstream Go project, then modified:
					//   https://github.com/golang/go/blob/3f995c3f3b43033013013e6c7ccc93a9b1411ca9/misc/wasm/wasm_exec.js#L404-L416
					"syscall/js.copyBytesToJS": (ret_addr, dest_addr, source_addr, source_len, source_cap) => {
						let num_bytes_copied_addr = ret_addr;
						let returned_status_addr = ret_addr + 4; // Address of returned boolean status variable

						const dst = loadValue(dest_addr);
						const src = loadSlice(source_addr, source_len);
						if (!(dst instanceof Uint8Array)) {
							mem().setUint8(returned_status_addr, 0); // Return "not ok" status
							return;
						}
						const toCopy = src.subarray(0, dst.length);
						dst.set(toCopy);
						setInt64(num_bytes_copied_addr, toCopy.length);
						mem().setUint8(returned_status_addr, 1); // Return "ok" status
					},
				}
			};
		}

		async run(instance) {
			this._inst = instance;
			this._values = [ // JS values that Go currently has references to, indexed by reference id
				NaN,
				0,
				null,
				true,
				false,
				global,
				this,
			];
			this._goRefCounts = []; // number of references that Go has to a JS value, indexed by reference id
			this._ids = new Map();  // mapping from JS values to reference ids
			this._idPool = [];      // unused ids that have been garbage collected
			this.exited = false;    // whether the Go program has exited

			const mem = new DataView(this._inst.exports.memory.buffer)

			while (true) {
				const callbackPromise = new Promise((resolve) => {
					this._resolveCallbackPromise = () => {
						if (this.exited) {
							throw new Error("bad callback: Go program has already exited");
						}
						setTimeout(resolve, 0); // make sure it is asynchronous
					};
				});
				this._inst.exports._start();
				if (this.exited) {
					break;
				}
				await callbackPromise;
			}
		}

		_resume() {
			if (this.exited) {
				throw new Error("Go program has already exited");
			}
			this._inst.exports.resume();
			if (this.exited) {
				this._resolveExitPromise();
			}
		}

		_makeFuncWrapper(id) {
			const go = this;
			return function () {
				const event = { id: id, this: this, args: arguments };
				go._pendingEvent = event;
				go._resume();
				return event.result;
			};
		}
	}

	

		const go = new Go();
		WebAssembly.instantiateStreaming(fetch("/assets/main.wasm"), go.importObject).then((result) => {
			return go.run(result.instance);
		}).catch((err) => {
			console.error(err);
		});
})();

			`),
			// Script(`const enosys = () => { const a = new Error("not implemented"); return a.code = "ENOSYS", a }; let outputBuf = ""; window.fs = { constants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, writeSync(a, b) { outputBuf += decoder.decode(b); const c = outputBuf.lastIndexOf("\n"); return -1 != c && (console.log(outputBuf.substr(0, c)), outputBuf = outputBuf.substr(c + 1)), b.length }, write(a, b, c, d, e, f) { if (0 !== c || d !== b.length || null !== e) return void f(enosys()); const g = this.writeSync(a, b); f(null, g) } }; const encoder = new TextEncoder("utf-8"), decoder = new TextDecoder("utf-8"); class Go { constructor() { this.argv = ["js"], this.env = {}, this.exit = a => { 0 !== a && console.warn("exit code:", a) }, this._exitPromise = new Promise(a => { this._resolveExitPromise = a }), this._pendingEvent = null, this._scheduledTimeouts = new Map, this._nextCallbackTimeoutID = 1; const a = (a, b) => { this.mem.setUint32(a + 0, b, !0), this.mem.setUint32(a + 4, Math.floor(b / 4294967296), !0) }, b = a => { const b = this.mem.getUint32(a + 0, !0), c = this.mem.getInt32(a + 4, !0); return b + 4294967296 * c }, c = a => { const b = this.mem.getFloat64(a, !0); if (0 !== b) { if (!isNaN(b)) return b; const c = this.mem.getUint32(a, !0); return this._values[c] } }, d = (a, b) => { const c = 2146959360; if ("number" == typeof b) return isNaN(b) ? (this.mem.setUint32(a + 4, 2146959360, !0), void this.mem.setUint32(a, 0, !0)) : 0 === b ? (this.mem.setUint32(a + 4, 2146959360, !0), void this.mem.setUint32(a, 1, !0)) : void this.mem.setFloat64(a, b, !0); switch (b) { case void 0: return void this.mem.setFloat64(a, 0, !0); case null: return this.mem.setUint32(a + 4, c, !0), void this.mem.setUint32(a, 2, !0); case !0: return this.mem.setUint32(a + 4, c, !0), void this.mem.setUint32(a, 3, !0); case !1: return this.mem.setUint32(a + 4, c, !0), void this.mem.setUint32(a, 4, !0); }let d = this._ids.get(b); d === void 0 && (d = this._idPool.pop(), d === void 0 && (d = this._values.length), this._values[d] = b, this._goRefCounts[d] = 0, this._ids.set(b, d)), this._goRefCounts[d]++; let e = 1; switch (typeof b) { case "string": e = 2; break; case "symbol": e = 3; break; case "function": e = 4; }this.mem.setUint32(a + 4, 2146959360 | e, !0), this.mem.setUint32(a, d, !0) }, e = a => { const c = b(a + 0), d = b(a + 8); return new Uint8Array(this._inst.exports.mem.buffer, c, d) }, f = d => { const e = b(d + 0), f = b(d + 8), g = Array(f); for (let a = 0; a < f; a++)g[a] = c(e + 8 * a); return g }, g = a => { const c = b(a + 0), d = b(a + 8); return decoder.decode(new DataView(this._inst.exports.mem.buffer, c, d)) }, h = Date.now() - performance.now(); this.importObject = { go: { "runtime.wasmExit": a => { const b = this.mem.getInt32(a + 8, !0); this.exited = !0, delete this._inst, delete this._values, delete this._goRefCounts, delete this._ids, delete this._idPool, this.exit(b) }, "runtime.wasmWrite": a => { const c = b(a + 8), d = b(a + 16), e = this.mem.getInt32(a + 24, !0); fs.writeSync(c, new Uint8Array(this._inst.exports.mem.buffer, d, e)) }, "runtime.resetMemoryDataView": () => { this.mem = new DataView(this._inst.exports.mem.buffer) }, "runtime.nanotime1": b => { a(b + 8, 1e6 * (h + performance.now())) }, "runtime.walltime1": b => { const c = new Date().getTime(); a(b + 8, c / 1e3), this.mem.setInt32(b + 16, 1e6 * (c % 1e3), !0) }, "runtime.scheduleTimeoutEvent": a => { const c = this._nextCallbackTimeoutID; this._nextCallbackTimeoutID++, this._scheduledTimeouts.set(c, setTimeout(() => { for (this._resume(); this._scheduledTimeouts.has(c);)console.warn("scheduleTimeoutEvent: missed timeout event"), this._resume() }, b(a + 8) + 1)), this.mem.setInt32(a + 16, c, !0) }, "runtime.clearTimeoutEvent": a => { const b = this.mem.getInt32(a + 8, !0); clearTimeout(this._scheduledTimeouts.get(b)), this._scheduledTimeouts.delete(b) }, "runtime.getRandomData": a => { crypto.getRandomValues(e(a + 8)) }, "syscall/js.finalizeRef": a => { const b = this.mem.getUint32(a + 8, !0); if (this._goRefCounts[b]--, 0 === this._goRefCounts[b]) { const a = this._values[b]; this._values[b] = null, this._ids.delete(a), this._idPool.push(b) } }, "syscall/js.stringVal": a => { d(a + 24, g(a + 8)) }, "syscall/js.valueGet": a => { const b = Reflect.get(c(a + 8), g(a + 16)); a = this._inst.exports.getsp(), d(a + 32, b) }, "syscall/js.valueSet": a => { Reflect.set(c(a + 8), g(a + 16), c(a + 32)) }, "syscall/js.valueDelete": a => { Reflect.deleteProperty(c(a + 8), g(a + 16)) }, "syscall/js.valueIndex": a => { d(a + 24, Reflect.get(c(a + 8), b(a + 16))) }, "syscall/js.valueSetIndex": a => { Reflect.set(c(a + 8), b(a + 16), c(a + 24)) }, "syscall/js.valueCall": a => { try { const b = c(a + 8), e = Reflect.get(b, g(a + 16)), h = f(a + 32), i = Reflect.apply(e, b, h); a = this._inst.exports.getsp(), d(a + 56, i), this.mem.setUint8(a + 64, 1) } catch (b) { d(a + 56, b), this.mem.setUint8(a + 64, 0) } }, "syscall/js.valueInvoke": a => { try { const b = c(a + 8), e = f(a + 16), g = Reflect.apply(b, void 0, e); a = this._inst.exports.getsp(), d(a + 40, g), this.mem.setUint8(a + 48, 1) } catch (b) { d(a + 40, b), this.mem.setUint8(a + 48, 0) } }, "syscall/js.valueNew": a => { try { const b = c(a + 8), e = f(a + 16), g = Reflect.construct(b, e); a = this._inst.exports.getsp(), d(a + 40, g), this.mem.setUint8(a + 48, 1) } catch (b) { d(a + 40, b), this.mem.setUint8(a + 48, 0) } }, "syscall/js.valueLength": b => { a(b + 16, parseInt(c(b + 8).length)) }, "syscall/js.valuePrepareString": b => { const e = encoder.encode(c(b + 8) + ""); d(b + 16, e), a(b + 24, e.length) }, "syscall/js.valueLoadString": a => { const b = c(a + 8); e(a + 16).set(b) }, "syscall/js.valueInstanceOf": a => { this.mem.setUint8(a + 24, c(a + 8) instanceof c(a + 16)) }, "syscall/js.copyBytesToGo": b => { const d = e(b + 8), f = c(b + 32); if (!(f instanceof Uint8Array)) return void this.mem.setUint8(b + 48, 0); const g = f.subarray(0, d.length); d.set(g), a(b + 40, g.length), this.mem.setUint8(b + 48, 1) }, "syscall/js.copyBytesToJS": b => { const d = c(b + 8), f = e(b + 16); if (!(d instanceof Uint8Array)) return void this.mem.setUint8(b + 48, 0); const g = f.subarray(0, d.length); d.set(g), a(b + 40, g.length), this.mem.setUint8(b + 48, 1) }, debug: a => { console.log(a) } } } } async run(a) { this._inst = a, this.mem = new DataView(this._inst.exports.mem.buffer), this._values = [NaN, 0, null, !0, !1, window, this], this._goRefCounts = [], this._ids = new Map, this._idPool = [], this.exited = !1; let b = 4096; const c = a => { const c = b, d = encoder.encode(a + "\0"); return new Uint8Array(this.mem.buffer, b, d.length).set(d), b += d.length, 0 != b % 8 && (b += 8 - b % 8), c }, d = this.argv.length, e = []; this.argv.forEach(a => { e.push(c(a)) }), e.push(0); const f = Object.keys(this.env).sort(); f.forEach(a => { e.push(c(a + "=" + this.env[a])) }), e.push(0); const g = b; e.forEach(a => { this.mem.setUint32(b, a, !0), this.mem.setUint32(b + 4, 0, !0), b += 8 }), this._inst.exports.run(d, g), this.exited && this._resolveExitPromise(), await this._exitPromise } _resume() { if (this.exited) throw new Error("Go program has already exited"); this._inst.exports.resume(), this.exited && this._resolveExitPromise() } _makeFuncWrapper(a) { const b = this; return function () { const c = { id: a, this: this, args: arguments }; return b._pendingEvent = c, b._resume(), c.result } } } const go = new Go; WebAssembly.instantiateStreaming(fetch("/assets/main.wasm"), go.importObject).then(a => go.run(a.instance)).catch(a => console.error("could not load wasm", a));`),
		),
		Body(ui),
	).Html(w)
}
