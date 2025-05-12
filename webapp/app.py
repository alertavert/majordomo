import logging
import argparse

import streamlit as st
from streamlit_option_menu import option_menu

from majordomo.api import (ask_assistant, Conversation, )
from majordomo import (list_projects, list_assistants, list_conversations,
                       get_conversation_from_name, create_conversation, )
from utils import setup_logger, get_logger
import constants


def render_conversation(conversation: Conversation):
    log = get_logger()
    cached_conv = st.session_state.conversations.get(conversation.id)
    log.debug(f"Cached conversation: {cached_conv.title if cached_conv else 'Not Found'}")
    if not cached_conv:
        return
    for message in cached_conv.messages:
        for role, text in message.items():
            with st.chat_message(role):
                st.write(text)


# Parse command line arguments
def parse_args():
    # Streamlit passes its own arguments; however, we can bypass it, by using `--` as a delimiter:
    # Run the app like this:
    #   streamlit run app.py -- --debug --server localhost:5050
    parser = argparse.ArgumentParser(description="Majordomo Code Assistant")
    parser.add_argument("--debug", action="store_true", help="Enable debug logging")
    parser.add_argument(
        "--server", default=f"{constants.SERVER}:{constants.PORT}",
        help="API server address (host:port)"
        )
    # Parse known args to avoid errors with Streamlit's own arguments
    args, _ = parser.parse_known_args()
    return args


# Initialize logger
@st.cache_resource
def init_logger(debug: bool) -> logging.Logger:
    log_level = logging.DEBUG if debug else logging.INFO
    setup_logger(log_level)
    return get_logger()


def main():
    st.set_page_config(page_title="Majordomo")
    args = parse_args()
    log = init_logger(args.debug)
    log.debug("Logging initialized")

    # TODO: use args.server to initialize the ApiServer class once it's done.
    log.debug(f"Connecting to API Server at: {args.server}")

    if "conversations" not in st.session_state:
        log.debug("Creating conversations cache")
        st.session_state.conversations = {}
    else:
        log.debug(f"Conversation cache initialized: {st.session_state.conversations}")

    header_cols = st.columns([1, 2])
    with header_cols[0]:
        st.image("https://img.icons8.com/ios-filled/50/000000/database.png", width=50)
    with header_cols[1]:
        st.title("Majordomo")
    with st.sidebar:
        st.title("Coding Assistant")
        selected = option_menu(
            menu_title="Navigation", options=["Ask Majordomo", "About"],
            icons=["cloud-upload", "info-circle"], menu_icon="cast", default_index=0, )

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
                    log.debug(f"Conversations for {active_project}: {conversations}")
                    selected_conversation = st.selectbox(
                        "Conversation", [c.title for c in conversations]
                        )
            with add_col:
                st.html(
                    "<span style='margin: 0px; padding: 0px; font-size:12px'>New "
                    "Conversation</span>"
                    )
                popover = st.popover("", icon=":material/add:")
                with popover:
                    selected_assistant = st.segmented_control(
                        "Assistants", [o.name for o in assistants], selection_mode="single"
                    )
                    conv_name = st.text_input("Name of the new conversation")
                    st.button(
                        "Create", on_click=create_conversation,
                        args=[conv_name, selected_assistant]
                    )
        with prompt_area:
            if selected_conversation:
                conv = get_conversation_from_name(selected_conversation, active_project)
                if conv:
                    if not conv.id:
                        log.error(f"conversations {conv.title} retrieved from API has no thread_id")
                    else:
                        # Retrieve all messages from the cache
                        # TODO: If it's not in the cache, we need to retrieve the messages from OpenAI API.
                        if conv.id in st.session_state.conversations:
                            conv = st.session_state.conversations[conv.id]
                        else:
                            st.session_state.conversations[conv.id] = conv
                else:
                    # This is a newly created conversation
                    conv = create_conversation(conv_name, selected_assistant)
                    log.debug(f"Conversation {conv} created")
                prompt = st.chat_input("Ask Majordomo")
                if prompt:
                    conv.messages.append({"USER": prompt})
                    try:
                        with st.spinner("Asking Majordomo..."):
                            response = ask_assistant(prompt, conv.assistant, conv.id)
                            conv.messages.append({"ASSISTANT": response.get("message")})
                            if not conv.id:
                                conv.id = response.get("thread_id")
                                log.debug(f"New conversation assigned ID: {conv.id}")
                            # Update the conversations' cache, including all messages
                            st.session_state.conversations[conv.id] = conv
                    except Exception as e:
                        st.error(f"An error occurred: {str(e)}")
        with chat_area:
            if selected_conversation:
                conv = get_conversation_from_name(selected_conversation, active_project)
                log.debug(f"Conversation for {selected_conversation}: {conv}")
                render_conversation(conv)

    elif selected == "About":
        st.header("About Majordomo")
        st.write(
            """
                    Our system uses specialized AI agents to:
                    1. 🔍 Analyze the Code content
                    3. 🎯 Enrich your Code with the most relevant information
                    4. 👥 Query the LLM with enriched data
                    5. 💡 Provide answers to detailed technical questions
                    """
            )


if __name__ == "__main__":
    main()
