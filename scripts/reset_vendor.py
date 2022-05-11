"""Script to update vendor dependencies"""
import os


def main():
    """Main function"""
    gopath = os.environ["GOPATH"]
    terraform_path = os.path.join(*[gopath, "src", "github.com", "zscaler",
                                    "terraform-provider-zpa"])
    os.chdir(terraform_path)
    os.system("git reset HEAD")
    os.system("git stash")
    os.system("rm -rf vendor")
    os.system("git checkout vendor")


if __name__ == "__main__":
    main()