# ------------------------------------------- MAKE GAME ---------------------------------------------- #

def ShowError(*msgs: str, code: int = 1):
    if len(msgs) > 0 and isinstance(msgs[-1], int):
        code = msgs[-1]
        msgs = msgs[:-1]

    msg = " ".join(f"{x}" for x in msgs)
    print(("[ ERROR ]: " if sys.platform == 'win32' else "\033[1;31m[ ERROR ]:\033[91m ") + str(msg))
    print("\033[0m")
    sys.exit(code)


try:
    import shutil
    import subprocess as sp
    import time
    import sys
    import os
    import pathlib
    from pathlib import Path
except ImportError as err:
    ShowError(f"There seems to be a problem, try again later: {err}")


###################### CONST ########################
_bin = Path("bin")
platform = __import__("platform").system()

GoNotInstaledExit = 2
BuildFailedExit = 3
IncorrectArgsExit = 4
ErrorExit = 5

os.chdir(os.path.dirname(os.path.abspath(__file__)))
os.chdir("..")

try:
    def arrowAnimaion(foo: bool):
        print("Checking dependencies...")
        cols = shutil.get_terminal_size().columns
        arrow = ">"
        steps = cols + len(arrow)
        delay = 1.8 / steps

        for i in range(steps):
            if i < cols:
                print("\r" + " " * cols, end="", flush=True)
                print("\r" + "-" * max(0, i - len(arrow)) + arrow[:min(len(arrow), cols - i)], end="", flush=True)
            time.sleep(delay)

        time.sleep(0.2)
        print("\r" + " " * cols, end="", flush=True)
        if foo:
            print("\r" + " " * ((cols - len("Done")) // 2) + "\033[1;32mDone\033[0m" if platform != "Windown" else "Done", flush=True)
        else:
            print("\r" + " " * ((cols - len("Error")) // 2) + "\033[1;31mError\033[0m" if platform != "Windown" else "Error", flush=True)

    def check():
        try:
            result = sp.run(["go", "version"], stdout=sp.PIPE, stderr=sp.PIPE, text=True)
            if result.returncode != 0:
                ShowError("Go is not installed or not in PATH", code=GoNotInstaledExit)
                return False

            VerInfo = result.stdout.strip()
            if not VerInfo.startswith("go version"):
                ShowError(f"Unexpected output: {VerInfo}", code=GoNotInstaledExit)
                return False

            version = VerInfo.split(" ")[2][2:]
            parts = version.split(".")
            n = int(parts[0])
            ver = int(parts[1])
            patch = int(parts[2]) if len(parts) > 2 else 0
            pass

            MinN, MinVer, MinPatch = 1, 21, 10

            if (n, ver, patch) >= (MinN, MinVer, MinPatch):
                return True
            else:
                ShowError(f"it looks like you don't have Go installed, or it's not in $PATH; properly install Go and try again (Error code: {VerInfo})", code=GoNotInstaledExit)
                pass
            pass
        except sp.CalledProcessError as e:
            ShowError(f"It looks like Go is not installed or is not in $PATH", code=GoNotInstaledExit)
            return False
        except PermissionError:
            ShowError("It seems that the script does not have permissions to skimp on the game. Check that you are running it with the correct permissions")
            return False
        except (FileNotFoundError, FileExistsError) as e:
            ShowError(f"It looks like Go is not installed or is not in $PATH", code=GoNotInstaledExit)
            return False
        except Exception as err:
            ShowError(f"There seems to be a problem, try again later: {err}")
            return False
        pass
    pass

    def checkAnimation(animation: bool = True):
        Check = check()
        if animation:
            arrowAnimaion(Check)
        if not Check:
            sys.exit(250)
        else:
            ansimsg = f'\033[32;1mAll dependencies are installed!\033[0m'
            msg = 'All dependencies are installed!'

            print(f"{ansimsg if platform == 'Linux' or platform == 'Darwin' else msg}")
        pass

    def loading(proc, txt = "Building bot"):
        d = ""
        while proc.poll() is None:
            print(f"\r{txt}{d:<3}", end="", flush=True)
            d += "."
            if len(d) > 3:
                d = ""
            time.sleep(0.15)
        print(f"\r{txt}... \033[32;1mDone!\033[0m", flush=True)
        pass
    pass

    # -------------------------------------------- FLAGS PARSE --------------------------------------- #

    args = sys.argv[1:]
    _os = {"windows", "linux", "darwin", "unix"}
    _arch = {"amd64", "386", "arm", "arm64", "ppc64", "ppc32"}
    os_ = "-"
    arch_ = "-"
    opt_ = 1
    for arg in args:
        if arg.startswith("--os="):
            os_ = arg.split("=", 1)[1].lower()
            if os_ not in _os:
                ShowError(f"Invalid OS specified: {os_}. Supported OS are: {', '.join(_os)}", code=IncorrectArgsExit)
                pass
            pass
        elif arg.startswith("--arch="):
            arch_ = arg.split("=", 1)[1].lower()
            if arch_ not in _arch:
                ShowError(
                    f"Invalid architecture specified: {arch_}. Supported architectures are: {', '.join(_arch)}", code=IncorrectArgsExit)

        elif arg.startswith("--opt="):
            opt = arg.split("=", 1)[1].lower()
            if opt == "1":
                continue
            elif opt == "2":
                opt_ = 2
            elif opt == "3":
                opt_ = 3
            else:
                ShowError(f"Invalid option specified: {opt}. Supported options are: 1, 2, 3", code=IncorrectArgsExit)
        pass
    pass

    env = os.environ.copy()
    if os_ != "-":
        env["GOOS"] = os_
    if arch_ != "-":
        env["GOARCH"] = arch_
    GC = ['go', 'build']

    run = False
    turbo = False
    if "-h" in args or "--help" in args:
        def f(clr, text):
            if sys.platform == 'win32':
                return text
            return f'\033[{clr}m{text}\033[0m'

        print(f('35;1', '--help/-h:') + ' displays help')
        print(f('94;1', '--clean/-c:') + ' Cleans and removes compiled parts')
        print(f('32;1', '--run/-r:') + ' Runs the program after compilation')
        print(f('33;1', "--static/-s:") + ' compiles with static linking')
        print(f('33;1', "--turbo/-t:") + 'Increases compilation speed without generating unnecessary graphic effects')
        print(f('92;1', "--opt={number}:") + "Uses optimizations to reduce the size of the final binaries. valid values: 1, 2, 3")
        print(f('35;1', "--os={os}:"), "Compiles the game for a given operating system")
        print(f('35;1', "--arch={arch}:"), "Compiles the game for a given procesor architecture")
        print(f('32;1', '--debug/-d'), "Compiles in debug mode")
        exit()
    if "-c" in args or "--clean" in args:
        if not _bin.exists():
            print("Already cleaned!" if sys.platform == 'win32' else "\033[32;1mAlready cleaned!\033[0m")
            exit()
        print("Cleaning...")
        try:
            if _bin.exists():
                print(f'deletion of {_bin.name}' if sys.platform == 'win32' else f'\033[91[deletion of {_bin.name}')
                shutil.rmtree(_bin)
            time.sleep(0.5)
            print('\033[0m' + '... Done' if sys.platform == 'win32' else '\033[32;1m... Done\033[0m')
        except (KeyboardInterrupt, SystemExit):
            exit()
        except (PermissionError, OSError):
            ShowError("There was an error while cleaning, check if the script is run with the right permissions")
        except:
            ShowError("There was an error while cleaning, try again")
        finally:
            exit()
    if '--static' in args or '-s' in args:
        print("\033[38;5;208m\033[1m[ WARNING ]\033[0m Static linking may not work if you do not have specific libraries installed, it is recommended to use dynamic linking (remove --static/-s flag)")
        GC.append('-ldflags="-extldflags \'-static\'"')
        env["CGO_ENABLED"] = "0"
    if "-r" in args or "--run" in args:
        run = True
    if '--turbo' in args or '-t' in args:
        turbo = True
    if opt_ == 2 or opt_ == 3:
        GC.append('-ldflags="-s -w"')
    if '--debug' in args or '-d' in args:
        GC.append('-tags')
        GC.append("DEBUG")


    checkAnimation(not turbo)
    print()

    # ---------------------------------------------- BUILD ------------------------------------------------ #
    try:
        cmd = sp.Popen([*GC, '-o', "../bin/PixelEconomy.exe" if (platform == "Windows" or os_ == "windows") else '../bin/PixelEconomy.elf'], cwd='src', env=env)
        loading(cmd)
        if cmd.returncode != 0:
            raise sp.CalledProcessError(cmd.returncode, cmd)
        if opt_ == 3:
            try:
                cmd = sp.Popen(['upx', '--best', '--lzma', 'bin/PixelEconomy.elf' if (platform == 'Windows' or os_ == 'windows') else 'bin/PixelEconomy.elf'], stdout=sp.DEVNULL, stderr=sp.DEVNULL)
                loading(cmd, "Performs optimizations on the binary")
                if cmd.returncode != 0:
                    raise sp.CalledProcessError(cmd.returncode, cmd)
            except sp.CalledProcessError:
                ShowError('It looks like the optymalization binary failed, try again later', code=ErrorExit)
            except PermissionError:
                ShowError('It seems that the script does not have permissions to compile the game. Check that you are running it with the correct permissions')
            except Exception as err:
                ShowError(f"There seems to be a problem, try again later: {err}", code=ErrorExit)
    except sp.CalledProcessError:
        ShowError('It looks like the compilation failed, try again later', code=BuildFailedExit)
    except (FileNotFoundError, OSError) as e:
        # DEBUG: print(e)
        ShowError(f"It looks like go is not installed or is not in $PATH", code=GoNotInstaledExit)
    except PermissionError:
        ShowError('It seems that the script does not have permissions to compile the game. Check that you are running it with the correct permissions')
    except Exception as err:
        ShowError(f"There seems to be a problem, try again later: {err}", code=BuildFailedExit)

    if run:
        try:
            sp.run('bin\\PixelEconomy.exe' if (platform == "Windows" or os_ == 'windows') else 'bin/PixelEconomy.elf', shell=True, check=True)
        except sp.CalledProcessError:
            ShowError(f'It looks like there was a problem while launching the game, try to launch it manually (found in {"bin/PixelEconomy.exe" if platform == 'Windows' else 'bin/PixelEconomy'})')
        except PermissionError:
            ShowError('It seems that the script does not have permissions to skimp on the game. Check that you are running it with the correct permissions')
        except Exception as err:
            ShowError(f"There seems to be a problem, try again later: {err}")
except KeyboardInterrupt:
    exit(0)
except Exception as err:
    ShowError(f"There seems to be a problem, try again later: {err}")
finally:
    print("\033[0m", end='')