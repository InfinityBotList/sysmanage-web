#!/usr/bin/env python3
import subprocess
import typing
import sys
import pathlib
import pydantic
import tqdm
import shutil
import os
import contextlib

class Step(pydantic.BaseModel):
    """Name of the current step"""
    name: str
    """Function to call to execute the step"""
    func: typing.Callable

class SingletonClass(object):
    def __new__(cls):
        if not hasattr(cls, 'instance'):
            cls.instance = super(SingletonClass, cls).__new__(cls)
        return cls.instance

class DummyFile(object):
    file = None
    def __init__(self, rds: "ProjectBuilder", file):
        self.file = file
        self.rds = rds

    def write(self, x):
        # Avoid print() second call (useless \n)
        if len(x.rstrip()) > 0:
            self.rds.pbar.write(x, file=self.file)
    
    def flush(self):
        pass

    def fileno(self):
        return self.file.fileno()

class ProjectBuilder(SingletonClass):
    repo_root: str
    target_dir: str
    step: typing.List[Step]
    pbar: tqdm.tqdm

    def __init__(
            self,
        ):
        self.repo_root = self.get_repo_root()
        self.target_dir = None
        self.pbar = None
        self.steps = []
    
    def get_input(self, s: str):
        """Get input from the user"""
        try:
            return input(s + ": ")
        except (KeyboardInterrupt, EOFError):
            print("Aborting due to KeyboardInterrupt or EOFError")
            sys.exit(1)
            

    def get_input_bool(self, s: str):
        """Get input from the user, but only allow y/n"""
        while True:
            inp = self.get_input(s + " (y/n)")
            if inp == "y":
                return True
            elif inp == "n":
                return False
            else:
                print("Invalid input, please enter y or n")
    
    def get_input_choice(self, s: str, choices: typing.List[str]):
        """Get input from the user, but only allow one of the choices"""
        while True:
            inp = self.get_input(s + " ({})".format(", ".join(choices)))
            if inp in choices:
                return inp
            else:
                print("Invalid input, please enter one of the choices")
    
    def get_input_int(self, s: str):
        """Get input from the user, but only allow integers"""
        while True:
            inp = self.get_input(s)
            try:
                return int(inp)
            except ValueError:
                print("Invalid input, please enter an integer")
    
    def get_repo_root(self):
        """Call git", "rev-parse", "--show-toplevel to get the root of the git repo"""
        return subprocess.check_output(["git", "rev-parse", "--show-toplevel"]).decode("utf-8").strip()

    @contextlib.contextmanager
    def _redirect(rds):
        save_stdout = sys.stdout
        sys.stdout = DummyFile(rds, sys.stdout)
        yield
        sys.stdout = save_stdout
    
    def exec(self, cmd: typing.List[str]):
        """Execute a command, redirecting stdout to tqdm.write and returning the output"""
        proj.pbar.write(f"> `{' '.join(cmd)}`")
        with contextlib.redirect_stdout(DummyFile(self, sys.stdout)):
            return subprocess.check_call(cmd, stdout=sys.stdout, stderr=sys.stdout)

    def main(self, target_dir: str):
        self.repo_root = self.get_repo_root()
        self.target_dir = target_dir

        print("Welcome to sysmanage project creator!")
        print("Repo root: {}".format(self.repo_root))
        print("Target dir: {}".format(self.target_dir))

        if self.repo_root == self.target_dir:
            print("Cannot create project in the root of the repo")
            sys.exit(1)
        
        if pathlib.Path(self.target_dir).exists():
            confirm = self.get_input_bool("Target dir already exists, continuing will *DELETE* this directory? (y/n)")
            
            if not confirm:
                print("Aborting")
                sys.exit(1)
            
            # Delete target dir
            subprocess.check_call(["rm", "-rf", self.target_dir])
        
        self.pbar = tqdm.tqdm(total=100, file=sys.stdout)

        with self._redirect():
            for step in self.steps:
                self.pbar.update(round(100 / (len(self.steps) + 1)))
                self.pbar.set_description(step.name)
                step.func()
            
            self.pbar.update(round(100 / (len(self.steps) + 1)))

            self.pbar.close()
    
    def step(self, name: str):
        """Decorator for adding steps"""
        def decorator(func):
            self.steps.append(Step(name=name, func=func))
            return func
        return decorator

proj = ProjectBuilder()

@proj.step(name="Create target dir")
def step_create_target_dir():
    pathlib.Path(proj.target_dir).mkdir(parents=True, exist_ok=True)

@proj.step(name="Copy template")
def step_copy_template():
    # Copy contents of example dir to target dir using shutil
    shutil.copytree(
        src=pathlib.Path(proj.repo_root) / "example",
        dst=pathlib.Path(proj.target_dir),
        dirs_exist_ok=True
    )

@proj.step(name="Setup project")
def step_setup_project():
    pt = pathlib.Path(proj.target_dir)
    os.chdir(pt)
    proj.exec(["go", "mod", "tidy"])
    proj.pbar.update(5)
    proj.exec(["git", "init"])
    proj.pbar.update(5)
    os.chdir("frontend")
    proj.exec(["npm", "install"])
    proj.pbar.update(5)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: {} <target_dir>".format(sys.argv[0]))
        sys.exit(1)

    proj.main(target_dir=sys.argv[1])
