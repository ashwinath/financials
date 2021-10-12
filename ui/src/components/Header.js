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

import { useLocation, useHistory } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { logoutAsync } from '../redux/loginSlice';
import { capitaliseFirstLetter } from '../utils';

export function Header() {
  const { loggedInUsername } = useSelector((state) => state.login)
  const history = useHistory();
  const { pathname } = useLocation();
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

  // Generate breadcrumbs
  const crumbs = pathname.split("/")
  crumbs[0] = "financials";
  const breadcrumbs = crumbs.map((item) => {
    return {
      text: capitaliseFirstLetter(item),
    };
  });

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
