
import streamlit as st

from majordomo.api import (ResponseError,
    get_projects,
    get_assistants,
    get_conversations,
)
from majordomo.model import Project, Assistant, Conversation


def show_error_response(resp: ResponseError) -> None:
    """
    Display an error response in the Streamlit app.
    """
    st.error(f"# {resp.title}\n{resp.message}", icon="⚠️")

@st.cache_data
def list_projects() -> tuple[str, list[Project]]:
    resp = get_projects()
    if isinstance(resp, ResponseError):
        show_error_response(resp)
        return "", []
    return resp


@st.cache_data
def list_assistants() -> list[Assistant]:
    resp = get_assistants()
    if isinstance(resp, ResponseError):
        show_error_response(resp)
        return []
    return resp


def list_conversations(project_name: str) -> list[Conversation]:
    resp = get_conversations(project_name)
    if isinstance(resp, ResponseError):
        show_error_response(resp)
        return []
    return resp


def get_project(name: str) -> Project:
    _, projects = list_projects()
    return next(p for p in projects if p.name == name)


def get_assistant(name: str) -> Assistant | None:
    for a in list_assistants():
        if a.name == name:
            return a
    else:
        show_error_response(ResponseError(
            title="Assistant Not Found",
            message=f"Assistant '{name}' does not exist in the system."
        ))
        return None


def get_conversation(conv_name: str, project_name: str) -> Conversation:
    return next(c for c in list_conversations(project_name) if c.title == conv_name)


def create_conversation(assistant: str, name: str | None = None) -> Conversation:
    return Conversation(
        id=None,
        title=name,
        assistant=assistant,
        messages=[],
    )
