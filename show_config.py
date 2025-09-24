import sys

def read_config(filename="config.txt"):
    config = {}
    with open(filename, "r") as f:
        for line in f:
            if "=" in line:
                key, value = line.strip().split("=", 1)
                config[key.strip()] = value.strip()
    return config

def write_config(config, filename="config.txt"):
    with open(filename, "w") as f:
        for key, value in config.items():
            f.write(f"{key} = {value}\n")

if __name__ == "__main__":
    config = read_config()

    # Parameter prüfen
    if len(sys.argv) == 3 and sys.argv[1] == "--set-account":
        new_account = sys.argv[2]
        config["account"] = new_account
        write_config(config)

    print("[core]")
    for key, value in config.items():
        print(f"{key} = {value}")
