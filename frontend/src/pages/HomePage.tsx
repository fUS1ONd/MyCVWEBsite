import React, { useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { profileService } from '../services';
import { Spinner, SEO } from '../components/common';

export const HomePage: React.FC = () => {
  const {
    data: profile,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['profile'],
    queryFn: profileService.getProfile,
  });

  // Smooth scroll for anchor links
  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      const href = target.getAttribute('href');
      if (href?.startsWith('#')) {
        e.preventDefault();
        const element = document.querySelector(href);
        element?.scrollIntoView({ behavior: 'smooth' });
      }
    };

    document.addEventListener('click', handleClick);
    return () => document.removeEventListener('click', handleClick);
  }, []);

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Ошибка загрузки</h2>
          <p className="text-gray-600">Не удалось загрузить информацию профиля</p>
        </div>
      </div>
    );
  }

  return (
    <>
      <SEO
        title="Главная"
        description={
          profile?.description || 'Персональный сайт с резюме и блогом об искусственном интеллекте'
        }
        image={profile?.photo_url}
        type="profile"
        keywords={['резюме', 'CV', 'портфолио', 'искусственный интеллект', 'AI', 'блог']}
      />
      <div className="bg-gray-50">
        {/* Hero Section */}
        <section id="hero" className="bg-gradient-to-br from-blue-50 to-indigo-100 py-20 md:py-32">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex flex-col md:flex-row items-center gap-12">
              {/* Photo */}
              <div className="flex-shrink-0">
                <div className="relative">
                  <div className="w-48 h-48 md:w-64 md:h-64 rounded-full overflow-hidden border-4 border-white shadow-xl">
                    {profile?.photo_url ? (
                      <img
                        src={profile.photo_url}
                        alt={profile.name}
                        className="w-full h-full object-cover"
                        loading="eager"
                      />
                    ) : (
                      <div className="w-full h-full bg-gradient-to-br from-blue-400 to-indigo-600 flex items-center justify-center">
                        <span className="text-6xl text-white font-bold">
                          {profile?.name?.charAt(0) || 'U'}
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Intro */}
              <div className="flex-1 text-center md:text-left">
                <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 mb-4 animate-fade-in">
                  {profile?.name || 'Имя не указано'}
                </h1>
                <p className="text-xl md:text-2xl text-gray-700 mb-6 animate-fade-in-delay-1">
                  {profile?.description || 'Описание не указано'}
                </p>
                <div className="flex flex-wrap gap-3 justify-center md:justify-start animate-fade-in-delay-2">
                  <a
                    href="#about"
                    className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-md"
                  >
                    Обо мне
                  </a>
                  <a
                    href="#contacts"
                    className="px-6 py-3 bg-white text-blue-600 rounded-lg hover:bg-gray-50 transition-colors shadow-md"
                  >
                    Контакты
                  </a>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* About Section */}
        <section id="about" className="py-20 bg-white">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-8 text-center">
                О себе
              </h2>
              <div className="prose prose-lg max-w-none text-gray-700 whitespace-pre-line">
                {profile?.description || (
                  <p className="text-gray-500 text-center italic">Информация не указана</p>
                )}
              </div>
            </div>
          </div>
        </section>

        {/* Activity Section */}
        <section id="activity" className="py-20 bg-gray-50">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-8 text-center">
                Текущая деятельность
              </h2>
              <div className="prose prose-lg max-w-none text-gray-700 whitespace-pre-line">
                {profile?.activity || (
                  <p className="text-gray-500 text-center italic">Информация не указана</p>
                )}
              </div>
            </div>
          </div>
        </section>

        {/* Contacts Section */}
        <section id="contacts" className="py-20 bg-white">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-8 text-center">
                Контакты
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {profile?.contacts?.email && (
                  <a
                    href={`mailto:${profile.contacts.email}`}
                    className="flex items-center gap-4 p-6 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                  >
                    <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg
                        className="w-6 h-6 text-white"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                        />
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">Email</p>
                      <p className="text-gray-600">{profile.contacts.email}</p>
                    </div>
                  </a>
                )}

                {profile?.contacts?.github && (
                  <a
                    href={profile.contacts.github}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-4 p-6 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                  >
                    <div className="w-12 h-12 bg-gray-900 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">GitHub</p>
                      <p className="text-gray-600">Профиль</p>
                    </div>
                  </a>
                )}

                {profile?.contacts?.linkedin && (
                  <a
                    href={profile.contacts.linkedin}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-4 p-6 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                  >
                    <div className="w-12 h-12 bg-blue-700 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">LinkedIn</p>
                      <p className="text-gray-600">Профиль</p>
                    </div>
                  </a>
                )}

                {profile?.contacts?.vk && (
                  <a
                    href={profile.contacts.vk}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-4 p-6 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                  >
                    <div className="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M15.07 2H8.93C3.33 2 2 3.33 2 8.93v6.14C2 20.67 3.33 22 8.93 22h6.14c5.6 0 6.93-1.33 6.93-6.93V8.93C22 3.33 20.67 2 15.07 2zm3.58 14.21h-1.72c-.47 0-.61-.38-1.45-1.23-.73-.73-1.06-.83-1.24-.83-.25 0-.33.08-.33.49v1.12c0 .3-.1.48-1.09.48-1.59 0-3.36-.97-4.6-2.77-1.87-2.58-2.38-4.53-2.38-4.93 0-.18.08-.35.49-.35h1.72c.36 0 .5.17.64.55.71 2.05 1.89 3.84 2.38 3.84.18 0 .27-.08.27-.54v-2.11c-.06-1.13-.66-1.23-.66-1.63 0-.14.12-.28.3-.28h2.7c.3 0 .42.16.42.51v2.85c0 .3.14.42.23.42.18 0 .36-.12.73-.49 1.13-1.27 1.94-3.23 1.94-3.23.1-.23.27-.45.68-.45h1.72c.52 0 .63.27.52.63-.19.85-2.28 3.48-2.28 3.48-.15.25-.21.36 0 .64.15.21.65.64 1 1.03.61.67 1.08 1.22 1.21 1.61.12.4-.09.6-.5.6z" />
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">ВКонтакте</p>
                      <p className="text-gray-600">Профиль</p>
                    </div>
                  </a>
                )}
              </div>
            </div>
          </div>
        </section>
      </div>
    </>
  );
};
