const localImages = [
  '/images/22c7a4b6ee22bbf989edf8784004d6bd.jpg',
  '/images/2d87e2fd00411848675681246a91d828.jpg',
  '/images/34020e9dcedf9dcea841c680a1532fc0.jpg',
  '/images/38363a9b1f9138b465c66dac5cbda87f.jpg',
  '/images/4062a26f1a95889cdb336f1a1cc12b09.jpg',
  '/images/52da1c940782c0d06cdab946473d2b12.jpg',
  '/images/5bf21f07534ccf756dd7a9cb8db0e46b.jpg',
  '/images/93b0383b24faca43b93ab83d8f241193.jpg',
  '/images/aa2d9eef4d75cc974684b5b4a6687de0.jpg',
  '/images/e7451a982ef5997d6b3e23c251e69d99.jpg',
  '/images/f5e8b6fd47e350e545037a4c179cf6ce.jpg',
  '/images/z7790002097095_8d2a2df7c816a038e6ff3e65ecfccb9d.jpg',
  '/images/z7790002110516_e3d8a16a1a912cdb1778c8515c890f48.jpg',
  '/images/z7791141669733_034b98dc8fe8cb9a09db2fec6769e4ca.jpg',
]

// Mock data for development — replace with API calls
export const mockComics = Array.from({ length: 20 }, (_, i) => ({
  id: i + 1,
  title: [
    'Naruto Shippuden', 'One Piece', 'Demon Slayer', 'Attack on Titan',
    'Jujutsu Kaisen', 'Bleach', 'Dragon Ball Z', 'My Hero Academia',
    'Tokyo Ghoul', 'Death Note', 'Vinland Saga', 'Berserk',
    'Chainsaw Man', 'Spy × Family', 'Blue Lock', 'Hunter x Hunter',
    'Fullmetal Alchemist', 'Vagabond', 'Kingdom', 'Monster',
  ][i],
  slug: [
    'naruto', 'one-piece', 'demon-slayer', 'attack-on-titan',
    'jujutsu-kaisen', 'bleach', 'dragon-ball-z', 'my-hero-academia',
    'tokyo-ghoul', 'death-note', 'vinland-saga', 'berserk',
    'chainsaw-man', 'spy-family', 'blue-lock', 'hunter-x-hunter',
    'fullmetal-alchemist', 'vagabond', 'kingdom', 'monster',
  ][i],
  type: 'comic' as const,
  coverUrl: localImages[i % localImages.length],
  latestChapter: (i * 15) + 1,
  updatedAt: `${(i % 24) + 1}h trước`,
  views: (i * 150000) + 50000,
  rating: +(7 + (i % 30) / 10).toFixed(1),
  status: i % 5 === 0 ? 'completed' as const : 'ongoing' as const,
  chapterCount: (i * 20) + 10,
  isHot: i < 6,
  isNew: i >= 14,
  genres: [['Action', 'Adventure'], ['Romance', 'Drama'], ['Fantasy', 'Magic'], ['Horror', 'Thriller']][i % 4],
  description: 'Một câu chuyện phiêu lưu kỳ thú về những anh hùng dũng cảm đối mặt với thử thách khốc liệt để bảo vệ thế giới và những người thân yêu.',
}))

export const mockNovels = Array.from({ length: 12 }, (_, i) => ({
  id: i + 100,
  title: ['Tru Tiên', 'Đấu Phá Thương Khung', 'Phàm Nhân Tu Tiên', 'Tiên Nghịch', 'Ngã Dục Phong Thiên', 'Huyền Hoàng', 'Long Vũ Cửu Thiên', 'Võng Du Chi Vương', 'Thần Mộ', 'Thiên Đạo Đồ Thư Quán', 'Tu La Võ Thần', 'Vạn Cổ Thần Đế'][i],
  slug: ['tru-tien', 'dau-pha-thuong-khung', 'pham-nhan-tu-tien', 'tien-nghich', 'nga-duc-phong-thien', 'huyen-hoang', 'long-vu-cuu-thien', 'vong-du-chi-vuong', 'than-mo', 'thien-dao-do-thu-quan', 'tu-la-vo-than', 'van-co-than-de'][i],
  type: 'novel' as const,
  coverUrl: localImages[(i + 5) % localImages.length],
  latestChapter: (i * 100) + 100,
  updatedAt: `${(i % 12) + 1}h trước`,
  views: (i * 250000) + 100000,
  rating: +(7.5 + (i % 25) / 10).toFixed(1),
  status: i % 4 === 0 ? 'completed' as const : 'ongoing' as const,
  chapterCount: (i * 150) + 100,
  isHot: i < 4,
  isNew: i >= 9,
  genres: [['Tiên hiệp', 'Tu tiên'], ['Ngôn tình', 'Drama'], ['Kiếm hiệp', 'Hành động'], ['Huyền huyễn', 'Phiêu lưu']][i % 4],
  description: 'Truyện kể về hành trình tu luyện và phiêu lưu của nhân vật chính trong một thế giới huyền ảo đầy bí ẩn và hiểm nguy.',
}))

export const mockSlides = mockComics.slice(0, 5).map(c => ({
  ...c,
  description: c.description,
}))

export const mockGenres = [
  { name: 'Action', slug: 'action', count: 324, icon: '⚔️' },
  { name: 'Romance', slug: 'romance', count: 287, icon: '💕' },
  { name: 'Fantasy', slug: 'fantasy', count: 412, icon: '🧙' },
  { name: 'Horror', slug: 'horror', count: 156, icon: '👻' },
  { name: 'Slice of Life', slug: 'slice-of-life', count: 198, icon: '🌸' },
  { name: 'Sci-Fi', slug: 'sci-fi', count: 143, icon: '🚀' },
  { name: 'Comedy', slug: 'comedy', count: 267, icon: '😂' },
  { name: 'Drama', slug: 'drama', count: 334, icon: '🎭' },
  { name: 'Tiên hiệp', slug: 'tien-hiep', count: 891, icon: '⛅' },
  { name: 'Kiếm hiệp', slug: 'kiem-hiep', count: 445, icon: '🗡️' },
  { name: 'Ngôn tình', slug: 'ngon-tinh', count: 623, icon: '💞' },
  { name: 'Huyền huyễn', slug: 'huyen-huyen', count: 734, icon: '🌀' },
]
