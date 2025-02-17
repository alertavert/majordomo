from typing import Tuple, List

import streamlit as st
from streamlit_option_menu import option_menu

from majordomo.session import (
    get_projects,
    get_assistants,
    get_conversations, Project, Assistant, new_conversation, Conversation,
)


@st.cache_data
def list_projects() -> List[Project]:
    return get_projects()


@st.cache_data
def list_assistants() -> List[Assistant]:
    return get_assistants("pr_1234")


def list_conversations(project_name):
    project_id = get_project_from_name(project_name).id
    return get_conversations(project_id)


def get_project_from_name(project_name: str) -> Project:
    return next(p for p in list_projects() if p.name == project_name)


def get_assistant_from_name(assistant_name: str) -> Assistant | None:
    for a in list_assistants():
        if a.name == assistant_name:
            return a

def get_conversation_from_name(conv_name: str, project_name: str) -> Conversation:
    return next(c for c in list_conversations(project_name) if c.name == conv_name)

def main():
    st.set_page_config(page_title="Majordomo")

    header_cols = st.columns([1, 2])
    with header_cols[0]:
        st.image("https://img.icons8.com/ios-filled/50/000000/database.png", width=50)
    with header_cols[1]:
        st.title("Majordomo")
    with st.sidebar:
        st.title("Coding Assistant")
        selected = option_menu(
            menu_title="Navigation",
            options=["Ask Majordomo", "About"],
            icons=["cloud-upload", "info-circle"],
            menu_icon="cast",
            default_index=0,
        )

    if selected == "Ask Majordomo":
        table_col1, table_col2 = st.columns(2)
        projects = list_projects()
        operators = list_assistants()
        with table_col1:
            active_project = st.selectbox("Active Project:", [p.name for p in projects])
        with table_col2:
            selected_conversation = None
            if active_project:
                conversations = list_conversations(active_project)
                if len(conversations) > 0:
                    selected_conversation = st.selectbox("Select Column", [c.name for c in conversations])
                else:
                    popover = st.popover("New Conversation")
                    with popover:
                        selected_assistant = st.segmented_control(
                            "",
                            [o.name for o in operators],
                            selection_mode="single"
                        )
                        conv_name = st.text_input("Name of the new conversation")
                        if conv_name and selected_assistant:
                            new_conv = new_conversation(
                                name=conv_name,
                                assistant=get_assistant_from_name(selected_assistant),
                                project=get_project_from_name(active_project),
                            )
                            conversations.append(new_conv)
        st.write(f"{[conv.name for conv in conversations]}")
        if selected_conversation:
            conv = get_conversation_from_name(selected_conversation, active_project)
            st.write(f"Asking {conv.assistant.name} about {conv.name}")
            prompt = st.text_input("", label_visibility="collapsed",
                                   placeholder="Ask Majordomo Anything About Coding")
            st.code(f"Asking {conv.assistant.name} about {prompt}")
    elif selected == "About":
        st.header("About Majordomo")
        st.write("""
        Our system uses specialized AI agents to:
        1. 🔍 Analyze the Code content
        3. 🎯 Enrich your Code with the most relevant information
        4. 👥 Query the LLM with enriched data
        5. 💡 Provide answers to detailed technical questions
        """)


if __name__ == "__main__":
    main()
