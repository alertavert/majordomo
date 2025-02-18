from time import sleep
from typing import Tuple, List

import streamlit as st
from streamlit_option_menu import option_menu

from majordomo.api import (
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


def create_conversation(name, project, assistant):
    new_conversation(
        name=name,
        assistant=get_assistant_from_name(assistant),
        project=get_project_from_name(project),
    )

def render_conversation():
    for message in st.session_state.conversation:
        for role, text in message.items():
            with st.chat_message(role):
                st.write(text)

def main():
    st.set_page_config(page_title="Majordomo")
    if "conversation" not in st.session_state:
        st.session_state.conversation = []

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
        header = st.container(border=True)
        chat_area = st.container(height=500)
        prompt_area = st.container()
        with header:
            proj_col, conv_col, add_col = st.columns([0.4, 0.3, 0.2])
            projects = list_projects()
            assistants = list_assistants()
            with proj_col:
                active_project = st.selectbox("Project", [p.name for p in projects])
            with conv_col:
                selected_conversation = None
                if active_project:
                    conversations = list_conversations(active_project)
                    selected_conversation = st.selectbox("Conversation", [c.name for c in conversations])
            with add_col:
                st.html("<span style='margin: 0px; padding: 0px; font-size:12px'>New Conversation</span>")
                popover = st.popover("", icon=":material/add:")
                with popover:
                    selected_assistant = st.segmented_control(
                        "Assistants",
                        [o.name for o in assistants],
                        selection_mode="single"
                    )
                    conv_name = st.text_input("Name of the new conversation")
                    st.button(
                        "Create",
                        on_click=create_conversation,
                        args=[conv_name, active_project, selected_assistant]
                    )
        with chat_area:
            render_conversation()
        with prompt_area:
            if selected_conversation:
                conv = get_conversation_from_name(selected_conversation, active_project)
                prompt = st.chat_input("Ask Majordomo")
                if prompt:
                    st.session_state.conversation.append({"USER": prompt})
                    try:
                        with st.spinner("Asking Majordomo..."):

                            sleep(2)
                            st.session_state.conversation.append({"ASSISTANT": "This is a mock response."})
                    except Exception as e:
                        st.error(f"An error occurred: {str(e)}")
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
