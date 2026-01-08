import { Popover } from 'antd';
import React, { FC, useEffect, useState } from 'react';
import { useRecoilValue } from 'recoil';
import sanitizeHtml from 'sanitize-html';
import Graphemer from 'graphemer';

import dynamic from 'next/dynamic';
import classNames from 'classnames';
import ContentEditable from './ContentEditable';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom } from '../../stores/ClientConfigStore';
import { MessageType } from '../../../interfaces/socket-events';
import styles from './ChatTextField.module.scss';

// Lazy loaded components

const EmojiPicker = dynamic(() => import('./EmojiPicker').then(mod => mod.EmojiPicker), {
	ssr: false,
});

const SendOutlined = dynamic(() => import('@ant-design/icons/SendOutlined'), {
	ssr: false,
});

const SmileOutlined = dynamic(() => import('@ant-design/icons/SmileOutlined'), {
	ssr: false,
});

export type ChatTextFieldProps = {
	defaultText?: string;
	enabled: boolean;
	focusInput: boolean;
};

const characterLimit = 300;
const maxNodeDepth = 10;
const graphemer = new Graphemer();

const normalizeText = (text: string) =>
	(text || '')
		.toLowerCase()
		.normalize('NFD')
		.replace(/[\u0300-\u036f]/g, '');

const prohibitedTerms = [
	'golpe',
	'golpes',
	'golpista',
	'golpistas',
	'golp',
	'g0lpe',
	'g0lpista',
	'scam',
	'scammer',
	'fraude',
	'fraude financeira',
	'esquema',
	'esquema ilegal',
	'piramide',
	'esquema de piramide',
	'pirâmide',
	'esquema de pirâmide',
	'estelionato',
	'crime',
	'criminoso',
	'quadrilha',
	'safado',
	'safados',
	'ladrão',
	'ladrao',
	'ladrões',
	'ladroes',
	'vagabundo',
	'vagabundos',
	'pilantra',
	'pilantragem',
	'canalha',
	'mau caráter',
	'mau-caráter',
	'bandido',
	'bandidagem',
	'charlatão',
	'charlatao',
	'charlatões',
	'picareta',
	'picaretagem',
	'sacar',
	'saque',
	'não consigo sacar',
	'nao consigo sacar',
	'nao consigo fazer saque',
	'saque travado',
	'saque bloqueado',
	'saque pendente',
	'dinheiro preso',
	'dinheiro bloqueado',
	'dinheiro sumiu',
	'cadê meu dinheiro',
	'cade meu dinheiro',
	'cadê o dinheiro',
	'cade o dinheiro',
	'perdi dinheiro',
	'perdi meu dinheiro',
	'não recebi',
	'nao recebi',
	'não caiu',
	'nao caiu',
	'saldo travado',
	'saldo bloqueado',
	'estou com problemas',
	'estou com problema',
	'problema',
	'problemas',
	'isso é problema',
	'isso e problema',
	'não funciona',
	'nao funciona',
	'não está funcionando',
	'nao esta funcionando',
	'não resolve',
	'nao resolve',
	'ninguém responde',
	'ninguem responde',
	'suporte não responde',
	'suporte nao responde',
	'ninguém ajuda',
	'ninguem ajuda',
	'mentira',
	'mentiroso',
	'mentirosos',
	'enganação',
	'enganacao',
	'enganar',
	'enganando',
	'ilusão',
	'ilusao',
	'ilusão financeira',
	'ilusao financeira',
	'promessa falsa',
	'falsas promessas',
	'propaganda enganosa',
	'furada',
	'fria',
	'isso é fria',
	'isso e fria',
	'perda de tempo',
	'não recomendo',
	'nao recomendo',
	'cuidado',
	'tomem cuidado',
	'alerta',
	'aviso',
	'denúncia',
	'denuncia',
	'vou denunciar',
	'denunciar',
	'reclame aqui',
	'procon',
	'polícia',
	'policia',
	'isso é golpe',
	'isso e golpe',
	'é golpe',
	'e golpe',
	'isso é fraude',
	'isso e fraude',
	'vocês são golpistas',
	'voces sao golpistas',
	'empresa golpista',
	'site golpista',
	'roubo',
	'roubaram meu dinheiro',
	'perdi tudo',
	'fui enganado',
	'fui enganada',
];

const prohibitedTermsNormalized = prohibitedTerms.map(normalizeText);

const getNodeTextContent = (node, depth) => {
	let text = '';

	if (depth > maxNodeDepth) return text;
	if (node === null) return text;

	switch (node.nodeType) {
		case Node.CDATA_SECTION_NODE: // unlikely
		case Node.TEXT_NODE: {
			text = node.nodeValue;
			break;
		}
		case Node.ELEMENT_NODE: {
			switch (node.tagName.toLowerCase()) {
				case 'img': {
					text = node.getAttribute('alt') || '';
					break;
				}
				case 'br': {
					text = '\n';
					break;
				}
				case 'strong':
				case 'b': {
					/* markdown representation of bold/strong */
					text = '**';
					for (let i = 0; i < node.childNodes.length; i += 1) {
						text += getNodeTextContent(node.childNodes[i], depth + 1);
					}
					text += '**';
					break;
				}
				case 'em':
				case 'i': {
					/* markdown representation of italic/emphasis */
					text = '*';
					for (let i = 0; i < node.childNodes.length; i += 1) {
						text += getNodeTextContent(node.childNodes[i], depth + 1);
					}
					text += '*';
					break;
				}
				case 'p': {
					text = '\n';
					for (let i = 0; i < node.childNodes.length; i += 1) {
						text += getNodeTextContent(node.childNodes[i], depth + 1);
					}
					break;
				}
				case 'a':
				case 'span':
				case 'div': {
					for (let i = 0; i < node.childNodes.length; i += 1) {
						text += getNodeTextContent(node.childNodes[i], depth + 1);
					}
					break;
				}
				/* nodes which should specifically not be parsed */
				case 'script':
				case 'style': {
					break;
				}
				default: {
					text = node.textContent;
					break;
				}
			}
			break;
		}
		default:
			break;
	}
	return text;
};

const getTextContent = node => {
	const text = getNodeTextContent(node, 0)
		.replace(/^\s+/, '') /* remove leading whitespace */
		.replace(/\s+$/, '') /* remove trailing whitespace */
		.replace(/\n([^\n])/g, '  \n$1') /* single line break to markdown break */
		.replace(
			/(?:http[s]?:\/\/.)?(?:www\.)?[-a-zA-Z0-9@%._+~#=]{2,256}\.[a-z]{2,6}\b(?:[-a-zA-Z0-9@:%_+.~#?&//=]*)/g,
			'*',
		);
	return text;
};

const containsProhibitedContent = (message: string) => {
	const normalizedMessage = normalizeText(message);
	return prohibitedTermsNormalized.some(term => normalizedMessage.includes(term));
};

export const ChatTextField: FC<ChatTextFieldProps> = ({ defaultText, enabled, focusInput }) => {
	const [characterCount, setCharacterCount] = useState(defaultText?.length);
	const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
	const [contentEditable, setContentEditable] = useState(null);
	const [customEmoji, setCustomEmoji] = useState([]);

	const onRootRef = el => {
		setContentEditable(el);
	};

	const getCharacterCount = () => {
		const message = getTextContent(contentEditable);
		return graphemer.countGraphemes(message);
	};

	const sendMessage = () => {
		if (!websocketService) {
			console.log('websocketService is not defined');
			return;
		}

		const message = getTextContent(contentEditable);
		if (containsProhibitedContent(message)) {
			contentEditable.innerHTML = '';
			setCharacterCount(0);
			return;
		}
		const containsLink = message.match(
			/(?:http[s]?:\/\/.)?(?:www\.)?[-a-zA-Z0-9@%._+~#=]{2,256}\.[a-z]{2,6}\b(?:[-a-zA-Z0-9@:%_+.~#?&//=]*)/g,
		);
		if (containsLink) {
			return;
		}
		const count = graphemer.countGraphemes(message);
		if (count === 0 || count > characterLimit) return;

		websocketService.send({ type: MessageType.CHAT, body: message });
		contentEditable.innerHTML = '';
	};

	const insertTextAtEnd = (textToInsert: string) => {
		contentEditable.innerHTML += textToInsert;
	};

	const onEmojiSelect = emoji => {
		if (emoji.native) {
			insertTextAtEnd(emoji.native);
		} else {
			// Custom emoji images
			const html = `<img src="${emoji.src}" alt=":${emoji.name}:" title=":${emoji.name}:" class="emoji" />`;
			insertTextAtEnd(html);
		}
	};

	const onKeyDown = (e: React.KeyboardEvent) => {
		if (e.key === 'Enter' && !(e.shiftKey || e.metaKey || e.ctrlKey || e.altKey)) {
			e.preventDefault();
			sendMessage();
		}
	};

	const onPaste = evt => {
		evt.preventDefault();

		const clip = evt.clipboardData;
		const { types } = clip;
		const contentTypes = ['text/html', 'text/plain'];

		let content;

		for (let i = 0; i < contentTypes.length; i += 1) {
			const contentType = contentTypes[i];

			if (types.includes(contentType)) {
				content = clip.getData(contentType);
				break;
			}
		}

		if (!content) return;

		const sanitized = sanitizeHtml(content, {
			allowedTags: ['b', 'i', 'em', 'strong', 'a', 'br', 'p', 'img'],
			allowedAttributes: {
				img: ['class', 'alt', 'title', 'src'],
			},
			allowedClasses: {
				img: ['emoji'],
			},
			transformTags: {
				h1: 'p',
				h2: 'p',
				h3: 'p',
			},
		});

		// MDN lists this as deprecated, but it's the only way to save this paste
		// into the browser's Undo buffer. Plus it handles all the selection
		// deletion, caret positioning, etc automaticaly.
		if (sanitized) document.execCommand('insertHTML', false, sanitized);
	};

	const handleChange = () => {
		const count = getCharacterCount();
		setCharacterCount(count);

		if (count === 0 && contentEditable.children.length === 1) {
			/* if we have a single <br> element added by the browser, remove. */
			if (contentEditable.children[0].tagName.toLowerCase() === 'br') {
				contentEditable.removeChild(contentEditable.children[0]);
			}
		}
	};

	// Focus the input when the component mounts.
	useEffect(() => {
		if (!focusInput) {
			return;
		}
		document.getElementById('chat-input-content-editable').focus({ preventScroll: true });
	}, []);

	const getCustomEmoji = async () => {
		try {
			const response = await fetch(`/api/emoji`);
			const emoji = await response.json();
			setCustomEmoji(emoji);

			emoji.forEach(e => {
				const preImg = document.createElement('link');
				preImg.href = e.url;
				preImg.rel = 'preload';
				preImg.as = 'image';
				document.head.appendChild(preImg);
			});
		} catch (e) {
			console.error('cannot fetch custom emoji', e);
		}
	};

	useEffect(() => {
		getCustomEmoji();
	}, []);

	return (
		<div id="chat-input" className={styles.root}>
			<div
				className={classNames(
					styles.inputWrap,
					characterCount > characterLimit && styles.maxCharacters,
				)}
			>
				<ContentEditable
					id="chat-input-content-editable"
					html={defaultText || ''}
					placeholder={enabled ? 'Send a message to chat' : 'Chat is disabled'}
					disabled={!enabled}
					onKeyDown={onKeyDown}
					onContentChange={handleChange}
					onPaste={onPaste}
					onRootRef={onRootRef}
					style={{ whiteSpace: 'pre-wrap', width: '100%' }}
					role="textbox"
					aria-label="Chat text input"
				/>
				{enabled && (
					<div style={{ display: 'flex', paddingLeft: '5px' }}>
						<Popover
							content={
								<div className={styles.emojiPickerContainer}>
									<EmojiPicker customEmoji={customEmoji} onEmojiSelect={onEmojiSelect} />
								</div>
							}
							trigger="click"
							placement="topRight"
						>
							<button
								type="button"
								aria-label="Emoji picker"
								id="owncast-emoji-picker-button"
								className={styles.emojiButton}
								title="Emoji picker button"
							>
								<SmileOutlined />
							</button>
						</Popover>
						<button
							type="button"
							aria-label="Send message"
							id="owncast-send-message-button"
							className={styles.sendButton}
							title="Send message Button"
							onClick={sendMessage}
						>
							<SendOutlined />
						</button>
					</div>
				)}
			</div>
		</div>
	);
};
