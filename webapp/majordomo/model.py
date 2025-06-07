from typing import Tuple, List

import streamlit as st

from majordomo.api import (
    Assistant,
    Conversation,
    Project,
    ResponseError,
    get_projects,
    get_assistants,
    get_conversations,
)

def show_error_response(resp: ResponseError):
    """
    Display an error response in the Streamlit app.
    """
    st.error(f"# {resp.title}\n{resp.message}", icon="⚠️")

@st.cache_data
def list_projects() -> Tuple[str, List[Project]]:
    resp = get_projects()
    if isinstance(resp, ResponseError):
        show_error_response(resp)
        return "", []
    return resp


@st.cache_data
def list_assistants() -> List[Assistant]:
    resp = get_assistants()
    if isinstance(resp, ResponseError):
        show_error_response(resp)
        return []
    return resp


def list_conversations(project_name) -> List[Conversation]:
    resp = get_conversations(project_name)
    if isinstance(resp, ResponseError):
        show_error_response(resp)
        return []
    return resp


def get_project_from_name(project_name: str) -> Project:
    _, projects = list_projects()
    return next(p for p in projects if p.name == project_name)


def get_assistant_from_name(assistant_name: str) -> Assistant | None:
    for a in list_assistants():
        if a.name == assistant_name:
            return a
    else:
        show_error_response(ResponseError(
            title="Assistant Not Found",
            message=f"Assistant '{assistant_name}' does not exist in the system."
        ))
        return None


def get_conversation_from_name(conv_name: str, project_name: str) -> Conversation:
    return next(c for c in list_conversations(project_name) if c.title == conv_name)


def create_conversation(name, assistant) -> Conversation:
    return Conversation(
        id=None,
        title=name,
        assistant=assistant,
        messages=[],
    )
