"""Generates Changelog.MD from git history"""
import subprocess


def main():
    """Main function"""
    ver_start = "1.0."
    out = subprocess.check_output("git log", shell=True)
    lines = out.split("\n")
    groups = []
    group = []
    for line in lines:
        if line.startswith("commit "):
            groups.append(group)
            group = [line]
        else:
            group.append(line.replace("*", " "))
    groups = groups[1:]
    for group in groups:
        if group[1].startswith("Merge: "):
            del group[1]
    version_last = len(groups) - 1
    for index in range(0, version_last + 1):
        group = groups[index]
        version = ver_start + str(version_last)
        date = group[2][8:]
        group_new = []
        for line in group:
            if line:
                group_new.append("*  " + line)

        version_last -= 1
        del group_new[2]
        del group_new[0]
        del group_new[0]

        group_new = ["", "## " + version + " (" + date + ")", "", "CHANGES", ""] + group_new
        groups[index] = group_new
    final = ""
    for group in groups:
        final += "\n".join(group)
    # print(json.dumps(groups, indent=2))
    with open("CHANGELOG.md", "w") as fileh:
        fileh.write(final)


if __name__ == "__main__":
    main()