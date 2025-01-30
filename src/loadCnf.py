import toml

DEFAULT = {
    "MoneyIcon": "💴",
    "WorkDelay": 30,
    "DatabasePath": "economy.db",
    "WorkEarningsMin": 50,
    "WorkEarningsMax": 200,
    "TopMessagesChannelID": None,       # if None is not sending messages
}


def loadCnf() -> tuple[str, int, str, int, int, any]:
    try: cnf = toml.load("config.toml")
    except (FileNotFoundError, toml.TomlDecodeError): cnf = {}

    try:
        return str(cnf.get("MoneyIcon", DEFAULT["MoneyIcon"])), int(cnf.get("WorkDelay", DEFAULT["WorkDelay"])), str(cnf.get("DatabasePath", DEFAULT["DatabasePath"])), int(cnf.get("WorkEarningsMin", DEFAULT["WorkEarningsMin"])), int(cnf.get("WorkEarningsMax", DEFAULT["WorkEarningsMax"])), cnf.get("TopMessagesChannelID", DEFAULT["TopMessagesChannelID"])
    except Exception as e:
        print("Invalid data in config.toml! using default")
        return DEFAULT["MoneyIcon"], DEFAULT["WorkDelay"], DEFAULT["DatabasePath"], DEFAULT["WorkEarningsMin"], DEFAULT["WorkEarningsMax"], DEFAULT["TopMessagesChannelID"]

def createDefault():
    with open('config.toml', 'w') as f:
        f.write(toml.dumps(DEFAULT))