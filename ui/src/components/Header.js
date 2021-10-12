import React, { useState } from 'react';

import {
  EuiHeader,
  EuiHeaderLogo,
  EuiHeaderSectionItemButton,
  EuiAvatar,
  EuiFlexGroup,
  EuiPopover,
  EuiFlexItem,
  EuiText,
  EuiSpacer,
  EuiLink,
  htmlIdGenerator,
} from '@elastic/eui';

import { useHistory } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { logoutAsync } from '../redux/loginSlice';

export function Header() {
  const { loggedInUsername } = useSelector((state) => state.login)
  const history = useHistory();
  const renderLogo = (
    <EuiHeaderLogo
      iconType="logoElastic"
      href="#"
      onClick={(e) => {
        e.preventDefault();
        history.push("/");
      }}
      aria-label="Go to home page"
    />
  );

  // TODO: deal with this
  const breadcrumbs = [
    {
      text: 'Management',
      href: '#',
      onClick: (e) => {
        e.preventDefault();
      },
    },
    {
      text: 'Users',
      href: '#',
      onClick: (e) => {
        e.preventDefault();
      },
    },
    {
      text: 'Create',
    },
  ];

  const loginButton = (
    <HeaderUserMenu username={loggedInUsername}/>
  );

  const sections = [
    {
      items: [renderLogo],
      borders: 'right',
      breadcrumbs: breadcrumbs,
      breadcrumbProps: {
        'aria-label': 'Header sections breadcrumbs',
      },
    },
    {
      items: [loginButton],
    },
  ];
  return (
    <EuiHeader
      sections={sections}
      postion="fixed"
    />
  );
};

function HeaderUserMenu({username}) {
  const dispatch = useDispatch();
  const id = htmlIdGenerator()();
  const [isOpen, setIsOpen] = useState(false);
  const onMenuButtonClick = () => {
    setIsOpen(!isOpen);
  };
  const closeMenu = () => {
    setIsOpen(false);
  };

  const usernameFormatted = username ? username.toUpperCase() : "?"

  const button = (
    <EuiHeaderSectionItemButton
      aria-controls={id}
      aria-expanded={isOpen}
      aria-haspopup="true"
      aria-label="Account menu"
      isDisabled={!username}
      onClick={onMenuButtonClick}
    >
      <EuiAvatar name={usernameFormatted} size="s" />
    </EuiHeaderSectionItemButton>
  );

  return (
    <EuiPopover
      id={id}
      button={button}
      isOpen={isOpen}
      anchorPosition="downRight"
      closePopover={closeMenu}
      panelPaddingSize="none"
    >
      <div style={{ width: 320 }}>
        <EuiFlexGroup
          gutterSize="m"
          className="euiHeaderProfile"
          responsive={false}
        >
          <EuiFlexItem grow={false}>
            <EuiAvatar name={usernameFormatted} size="xl" />
          </EuiFlexItem>
          <EuiFlexItem>
            <EuiText>
              <p>{username}</p>
            </EuiText>
            <EuiSpacer size="m" />
            <EuiFlexGroup>
              <EuiFlexItem>
                <EuiFlexGroup justifyContent="spaceBetween">
                  <EuiFlexItem grow={false}>
                    <EuiLink onClick={() => dispatch(logoutAsync())}>Log out</EuiLink>
                  </EuiFlexItem>
                </EuiFlexGroup>
              </EuiFlexItem>
            </EuiFlexGroup>
          </EuiFlexItem>
        </EuiFlexGroup>
      </div>
    </EuiPopover>
  );
};
