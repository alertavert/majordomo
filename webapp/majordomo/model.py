from typing import Tuple, List

import streamlit as st

from majordomo.api import Project, get_projects, ResponseError, Assistant, get_assistants, get_conversations, \
    Conversation, new_conversation


@st.cache_data
def list_projects() -> Tuple[str, List[Project]]:
    resp = get_projects()
    if isinstance(resp, ResponseError):
        st.error(f"# {resp.title}\n{resp.message}", icon="⚠️")
        return "", []
    return resp


@st.cache_data
def list_assistants() -> List[Assistant]:
    resp = get_assistants()
    if isinstance(resp, ResponseError):
        st.error(f"# {resp.title}\n{resp.message}", icon="⚠️")
        return []
    return resp


@st.cache_data
def list_conversations(project_name):
    project_id = get_project_from_name(project_name).name
    resp = get_conversations(project_id)
    if isinstance(resp, ResponseError):
        st.error(f"# {resp.title}\n{resp.message}", icon="⚠️")
        return []
    return resp


def get_project_from_name(project_name: str) -> Project:
    _, projects = list_projects()
    return next(p for p in projects if p.name == project_name)


def get_assistant_from_name(assistant_name: str) -> Assistant | None:
    for a in list_assistants():
        if a.name == assistant_name:
            return a


def get_conversation_from_name(conv_name: str, project_name: str) -> Conversation:
    return next(c for c in list_conversations(project_name) if c.name == conv_name)


def create_conversation(name, project, assistant):
    new_conversation(
        name=name,
        assistant=get_assistant_from_name(assistant),
        project=get_project_from_name(project),
    )
