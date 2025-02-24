import logging
import sys

import streamlit as st
from streamlit_option_menu import option_menu

from majordomo.api import (
    ask_assistant,
)
from majordomo import (
    list_projects,
    list_assistants,
    list_conversations,
    get_conversation_from_name,
    create_conversation,
)
from utils import setup_logger, get_logger


def render_conversation():
    for message in st.session_state.conversation:
        for role, text in message.items():
            with st.chat_message(role):
                st.write(text)

# Initialize logger
@st.cache_resource
def init_logger() -> logging.Logger:
    log_level = logging.DEBUG if len(sys.argv) == 2 and sys.argv[1] == "debug" else logging.INFO
    setup_logger(log_level)
    return get_logger()

def main():
    st.set_page_config(page_title="Majordomo")
    logger = init_logger()

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
            # TODO: use active project to pre-select in the selectbox
            _, projects = list_projects()
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
                            response = ask_assistant(prompt, conv.id, conv.assistant)
                            st.session_state.conversation.append({"ASSISTANT": response})
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
