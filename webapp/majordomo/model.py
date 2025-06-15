from dataclasses import dataclass
from pathlib import Path
from typing import Any


@dataclass
class Project:
    name: str
    description: str
    location: Path

    @staticmethod
    def from_dict(data: dict[str, Any]) -> "Project":
        return Project(
            name=data.get("name"), description=data.get("description"),
            location=Path(data.get("location"))
        )


@dataclass
class Assistant:
    id: str
    name: str
    model: str
    instructions: str

    @staticmethod
    def from_dict(data: dict[str, Any]) -> "Assistant":
        return Assistant(
            id=data.get("id"), name=data.get("name"), model=data.get("model"),
            instructions=data.get("instructions")
        )


@dataclass
class Conversation:
    id: str | None
    title: str
    assistant: str
    messages: list[dict[str, str]]

    @staticmethod
    def from_dict(data: dict[str, Any]) -> "Conversation":
        return Conversation(
            id=data.get("id"),
            title=data.get("name"),
            assistant=data.get("assistant"),
            messages=[]
        )
