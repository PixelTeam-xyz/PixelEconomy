import nextcord

def genErr(*msg: str, title: str = "Błąd") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"❌ {title}",
        description=" ".join(msg),
        color=nextcord.Color.red(),
    )

def genInfo(*msg: str, title: str = "Informacja") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"ℹ️ {title}",
        description=" ".join(msg),
        color=nextcord.Color.blue(),
    )

def genTip(*msg: str, title: str = "Porada") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"💡 {title}",
        description=" ".join(msg),
        color=nextcord.Color.green(),
    )

def genWarn(*msg: str, title: str = "Ostrzeżenie") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"⚠️ {title}",
        description=" ".join(msg),
        color=nextcord.Color.gold(),
    )

def genSuccess(*msg: str, title: str = "Sukces") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"✅ {title}",
        description=" ".join(msg),
        color=nextcord.Color.green(),
    )

def genQuestion(*msg: str, title: str = "Pytanie") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"❓ {title}",
        description=" ".join(msg),
        color=nextcord.Color.blurple(),
    )

def genCustom(title: str, *msg: str, color: nextcord.Color = nextcord.Color.default(), emoji: str = "") -> nextcord.Embed:
    return nextcord.Embed(
        title=f"{emoji} {title}",
        description=" ".join(msg),
        color=color,
    )