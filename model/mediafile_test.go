package model_test

import (
	"time"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/conf/configtest"
	. "github.com/navidrome/navidrome/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MediaFiles", func() {
	var mfs MediaFiles

	Describe("ToAlbum", func() {
		Context("Simple attributes", func() {
			BeforeEach(func() {
				mfs = MediaFiles{
					{
						ID: "1", AlbumID: "AlbumID", Album: "Album", ArtistID: "ArtistID", Artist: "Artist", AlbumArtistID: "AlbumArtistID", AlbumArtist: "AlbumArtist",
						SortAlbumName: "SortAlbumName", SortArtistName: "SortArtistName", SortAlbumArtistName: "SortAlbumArtistName",
						OrderAlbumName: "OrderAlbumName", OrderAlbumArtistName: "OrderAlbumArtistName",
						MbzAlbumArtistID: "MbzAlbumArtistID", MbzAlbumType: "MbzAlbumType", MbzAlbumComment: "MbzAlbumComment",
						MbzReleaseGroupID: "MbzReleaseGroupID", Compilation: false, CatalogNum: "", Path: "/music1/file1.mp3", FolderID: "Folder1",
					},
					{
						ID: "2", Album: "Album", ArtistID: "ArtistID", Artist: "Artist", AlbumArtistID: "AlbumArtistID", AlbumArtist: "AlbumArtist", AlbumID: "AlbumID",
						SortAlbumName: "SortAlbumName", SortArtistName: "SortArtistName", SortAlbumArtistName: "SortAlbumArtistName",
						OrderAlbumName: "OrderAlbumName", OrderArtistName: "OrderArtistName", OrderAlbumArtistName: "OrderAlbumArtistName",
						MbzAlbumArtistID: "MbzAlbumArtistID", MbzAlbumType: "MbzAlbumType", MbzAlbumComment: "MbzAlbumComment",
						MbzReleaseGroupID: "MbzReleaseGroupID",
						Compilation:       true, CatalogNum: "CatalogNum", HasCoverArt: true, Path: "/music2/file2.mp3", FolderID: "Folder2",
					},
				}
			})

			It("sets the single values correctly", func() {
				album := mfs.ToAlbum()
				Expect(album.ID).To(Equal("AlbumID"))
				Expect(album.Name).To(Equal("Album"))
				Expect(album.AlbumArtist).To(Equal("AlbumArtist"))
				Expect(album.AlbumArtistID).To(Equal("AlbumArtistID"))
				Expect(album.SortAlbumName).To(Equal("SortAlbumName"))
				Expect(album.SortAlbumArtistName).To(Equal("SortAlbumArtistName"))
				Expect(album.OrderAlbumName).To(Equal("OrderAlbumName"))
				Expect(album.OrderAlbumArtistName).To(Equal("OrderAlbumArtistName"))
				Expect(album.MbzAlbumArtistID).To(Equal("MbzAlbumArtistID"))
				Expect(album.MbzAlbumType).To(Equal("MbzAlbumType"))
				Expect(album.MbzAlbumComment).To(Equal("MbzAlbumComment"))
				Expect(album.MbzReleaseGroupID).To(Equal("MbzReleaseGroupID"))
				Expect(album.CatalogNum).To(Equal("CatalogNum"))
				Expect(album.Compilation).To(BeTrue())
				Expect(album.EmbedArtPath).To(Equal("/music2/file2.mp3"))
				Expect(album.FolderIDs).To(ConsistOf("Folder1", "Folder2"))
			})
		})
		Context("Aggregated attributes", func() {
			When("we don't have any songs", func() {
				BeforeEach(func() {
					mfs = MediaFiles{}
				})
				It("returns an empty album", func() {
					album := mfs.ToAlbum()
					Expect(album.Duration).To(Equal(float32(0)))
					Expect(album.Size).To(Equal(int64(0)))
					Expect(album.MinYear).To(Equal(0))
					Expect(album.MaxYear).To(Equal(0))
					Expect(album.Date).To(BeEmpty())
					Expect(album.UpdatedAt).To(BeZero())
					Expect(album.CreatedAt).To(BeZero())
				})
			})
			When("we have only one song", func() {
				BeforeEach(func() {
					mfs = MediaFiles{
						{Duration: 100.2, Size: 1024, Year: 1985, Date: "1985-01-02", UpdatedAt: t("2022-12-19 09:30"), BirthTime: t("2022-12-19 08:30")},
					}
				})
				It("calculates the aggregates correctly", func() {
					album := mfs.ToAlbum()
					Expect(album.Duration).To(Equal(float32(100.2)))
					Expect(album.Size).To(Equal(int64(1024)))
					Expect(album.MinYear).To(Equal(1985))
					Expect(album.MaxYear).To(Equal(1985))
					Expect(album.Date).To(Equal("1985-01-02"))
					Expect(album.UpdatedAt).To(Equal(t("2022-12-19 09:30")))
					Expect(album.CreatedAt).To(Equal(t("2022-12-19 08:30")))
				})
			})

			When("we have multiple songs with different dates", func() {
				BeforeEach(func() {
					mfs = MediaFiles{
						{Duration: 100.2, Size: 1024, Year: 1985, Date: "1985-01-02", UpdatedAt: t("2022-12-19 09:30"), BirthTime: t("2022-12-19 08:30")},
						{Duration: 200.2, Size: 2048, Year: 0, Date: "", UpdatedAt: t("2022-12-19 09:45"), BirthTime: t("2022-12-19 08:30")},
						{Duration: 150.6, Size: 1000, Year: 1986, Date: "1986-01-02", UpdatedAt: t("2022-12-19 09:45"), BirthTime: t("2022-12-19 07:30")},
					}
				})
				It("calculates the aggregates correctly", func() {
					album := mfs.ToAlbum()
					Expect(album.Duration).To(Equal(float32(451.0)))
					Expect(album.Size).To(Equal(int64(4072)))
					Expect(album.MinYear).To(Equal(1985))
					Expect(album.MaxYear).To(Equal(1986))
					Expect(album.Date).To(BeEmpty())
					Expect(album.UpdatedAt).To(Equal(t("2022-12-19 09:45")))
					Expect(album.CreatedAt).To(Equal(t("2022-12-19 07:30")))
				})
				Context("MinYear", func() {
					It("returns 0 when all values are 0", func() {
						mfs = MediaFiles{{Year: 0}, {Year: 0}, {Year: 0}}
						a := mfs.ToAlbum()
						Expect(a.MinYear).To(Equal(0))
					})
					It("returns the smallest value from the list, not counting 0", func() {
						mfs = MediaFiles{{Year: 2000}, {Year: 0}, {Year: 1999}}
						a := mfs.ToAlbum()
						Expect(a.MinYear).To(Equal(1999))
					})
				})
			})
			When("we have multiple songs with same dates", func() {
				BeforeEach(func() {
					mfs = MediaFiles{
						{Duration: 100.2, Size: 1024, Year: 1985, Date: "1985-01-02", UpdatedAt: t("2022-12-19 09:30"), BirthTime: t("2022-12-19 08:30")},
						{Duration: 200.2, Size: 2048, Year: 1985, Date: "1985-01-02", UpdatedAt: t("2022-12-19 09:45"), BirthTime: t("2022-12-19 08:30")},
						{Duration: 150.6, Size: 1000, Year: 1985, Date: "1985-01-02", UpdatedAt: t("2022-12-19 09:45"), BirthTime: t("2022-12-19 07:30")},
					}
				})
				It("sets the date field correctly", func() {
					album := mfs.ToAlbum()
					Expect(album.Date).To(Equal("1985-01-02"))
					Expect(album.MinYear).To(Equal(1985))
					Expect(album.MaxYear).To(Equal(1985))
				})
			})
			DescribeTable("explicitStatus",
				func(mfs MediaFiles, status string) {
					Expect(mfs.ToAlbum().ExplicitStatus).To(Equal(status))
				},
				Entry("sets the album to clean when a clean song is present", MediaFiles{{ExplicitStatus: ""}, {ExplicitStatus: "c"}, {ExplicitStatus: ""}}, "c"),
				Entry("sets the album to explicit when an explicit song is present", MediaFiles{{ExplicitStatus: ""}, {ExplicitStatus: "e"}, {ExplicitStatus: ""}}, "e"),
				Entry("takes precedence of explicit songs over clean ones", MediaFiles{{ExplicitStatus: "e"}, {ExplicitStatus: "c"}, {ExplicitStatus: ""}}, "e"),
			)
		})
		Context("Calculated attributes", func() {
			Context("Discs", func() {
				When("we have no discs info", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{Album: "Album1"}, {Album: "Album1"}, {Album: "Album1"}}
					})
					It("adds 1 disc without subtitle", func() {
						album := mfs.ToAlbum()
						Expect(album.Discs).To(Equal(Discs{1: ""}))
					})
				})
				When("we have only one disc", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{DiscNumber: 1, DiscSubtitle: "DiscSubtitle"}}
					})
					It("sets the correct Discs", func() {
						album := mfs.ToAlbum()
						Expect(album.Discs).To(Equal(Discs{1: "DiscSubtitle"}))
					})
				})
				When("we have multiple discs", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{DiscNumber: 1, DiscSubtitle: "DiscSubtitle"}, {DiscNumber: 2, DiscSubtitle: "DiscSubtitle2"}, {DiscNumber: 1, DiscSubtitle: "DiscSubtitle"}}
					})
					It("sets the correct Discs", func() {
						album := mfs.ToAlbum()
						Expect(album.Discs).To(Equal(Discs{1: "DiscSubtitle", 2: "DiscSubtitle2"}))
					})
				})
			})

			Context("Genres/tags", func() {
				When("we don't have any tags", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{}}
					})
					It("sets the correct Genre", func() {
						album := mfs.ToAlbum()
						Expect(album.Tags).To(BeEmpty())
					})
				})
				When("we have only one Genre", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{Tags: Tags{"genre": []string{"Rock"}}}}
					})
					It("sets the correct Genre", func() {
						album := mfs.ToAlbum()
						Expect(album.Tags).To(HaveLen(1))
						Expect(album.Tags).To(HaveKeyWithValue(TagGenre, []string{"Rock"}))
					})
				})
				When("we have multiple Genres", func() {
					BeforeEach(func() {
						mfs = MediaFiles{
							{Tags: Tags{"genre": []string{"Punk"}, "mood": []string{"Happy", "Chill"}}},
							{Tags: Tags{"genre": []string{"Rock"}}},
							{Tags: Tags{"genre": []string{"Alternative", "Rock"}}},
						}
					})
					It("sets the correct Genre, sorted by frequency, then alphabetically", func() {
						album := mfs.ToAlbum()
						Expect(album.Tags).To(HaveLen(2))
						Expect(album.Tags).To(HaveKeyWithValue(TagGenre, []string{"Rock", "Alternative", "Punk"}))
						Expect(album.Tags).To(HaveKeyWithValue(TagMood, []string{"Chill", "Happy"}))
					})
				})
				When("we have tags with mismatching case", func() {
					BeforeEach(func() {
						mfs = MediaFiles{
							{Tags: Tags{"genre": []string{"synthwave"}}},
							{Tags: Tags{"genre": []string{"Synthwave"}}},
						}
					})
					It("normalizes the tags in just one", func() {
						album := mfs.ToAlbum()
						Expect(album.Tags).To(HaveLen(1))
						Expect(album.Tags).To(HaveKeyWithValue(TagGenre, []string{"Synthwave"}))
					})
				})
			})
			Context("Comments", func() {
				When("we have only one Comment", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{Comment: "comment1"}}
					})
					It("sets the correct Comment", func() {
						album := mfs.ToAlbum()
						Expect(album.Comment).To(Equal("comment1"))
					})
				})
				When("we have multiple equal comments", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{Comment: "comment1"}, {Comment: "comment1"}, {Comment: "comment1"}}
					})
					It("sets the correct Comment", func() {
						album := mfs.ToAlbum()
						Expect(album.Comment).To(Equal("comment1"))
					})
				})
				When("we have different comments", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{Comment: "comment1"}, {Comment: "not the same"}, {Comment: "comment1"}}
					})
					It("sets the correct comment", func() {
						album := mfs.ToAlbum()
						Expect(album.Comment).To(BeEmpty())
					})
				})
			})
			Context("Participants", func() {
				var album Album
				BeforeEach(func() {
					mfs = MediaFiles{
						{
							Album: "Album1", AlbumArtistID: "AA1", AlbumArtist: "Display AlbumArtist1", Artist: "Artist1",
							DiscSubtitle: "DiscSubtitle1", SortAlbumName: "SortAlbumName1",
							Participants: Participants{
								RoleAlbumArtist: ParticipantList{_p("AA1", "AlbumArtist1", "SortAlbumArtistName1")},
								RoleArtist:      ParticipantList{_p("A1", "Artist1", "SortArtistName1")},
							},
						},
						{
							Album: "Album1", AlbumArtistID: "AA1", AlbumArtist: "Display AlbumArtist1", Artist: "Artist2",
							DiscSubtitle: "DiscSubtitle2", SortAlbumName: "SortAlbumName1",
							Participants: Participants{
								RoleAlbumArtist: ParticipantList{_p("AA1", "AlbumArtist1", "SortAlbumArtistName1")},
								RoleArtist:      ParticipantList{_p("A2", "Artist2", "SortArtistName2")},
								RoleComposer:    ParticipantList{_p("C1", "Composer1")},
							},
						},
					}
					album = mfs.ToAlbum()
				})
				It("gets all participants from all tracks", func() {
					Expect(album.Participants).To(HaveKeyWithValue(RoleAlbumArtist, ParticipantList{_p("AA1", "AlbumArtist1", "SortAlbumArtistName1")}))
					Expect(album.Participants).To(HaveKeyWithValue(RoleComposer, ParticipantList{_p("C1", "Composer1")}))
					Expect(album.Participants).To(HaveKeyWithValue(RoleArtist, ParticipantList{
						_p("A1", "Artist1", "SortArtistName1"), _p("A2", "Artist2", "SortArtistName2"),
					}))
				})
			})
			Context("MbzAlbumID", func() {
				When("we have only one MbzAlbumID", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{MbzAlbumID: "id1"}}
					})
					It("sets the correct MbzAlbumID", func() {
						album := mfs.ToAlbum()
						Expect(album.MbzAlbumID).To(Equal("id1"))
					})
				})
				When("we have multiple MbzAlbumID", func() {
					BeforeEach(func() {
						mfs = MediaFiles{{MbzAlbumID: "id1"}, {MbzAlbumID: "id2"}, {MbzAlbumID: "id1"}}
					})
					It("uses the most frequent MbzAlbumID", func() {
						album := mfs.ToAlbum()
						Expect(album.MbzAlbumID).To(Equal("id1"))
					})
				})
			})
			Context("Album Art", func() {
				When("we have media files with cover art from multiple discs", func() {
					BeforeEach(func() {
						mfs = MediaFiles{
							{
								Path:        "Artist/Album/Disc2/01.mp3",
								HasCoverArt: true,
								DiscNumber:  2,
							},
							{
								Path:        "Artist/Album/Disc1/01.mp3",
								HasCoverArt: true,
								DiscNumber:  1,
							},
							{
								Path:        "Artist/Album/Disc3/01.mp3",
								HasCoverArt: true,
								DiscNumber:  3,
							},
						}
					})
					It("selects the cover art from the lowest disc number", func() {
						album := mfs.ToAlbum()
						Expect(album.EmbedArtPath).To(Equal("Artist/Album/Disc1/01.mp3"))
					})
				})

				When("we have media files with cover art from the same disc number", func() {
					BeforeEach(func() {
						mfs = MediaFiles{
							{
								Path:        "Artist/Album/Disc1/02.mp3",
								HasCoverArt: true,
								DiscNumber:  1,
							},
							{
								Path:        "Artist/Album/Disc1/01.mp3",
								HasCoverArt: true,
								DiscNumber:  1,
							},
						}
					})
					It("selects the cover art with the lowest path alphabetically", func() {
						album := mfs.ToAlbum()
						Expect(album.EmbedArtPath).To(Equal("Artist/Album/Disc1/01.mp3"))
					})
				})

				When("we have media files with some missing cover art", func() {
					BeforeEach(func() {
						mfs = MediaFiles{
							{
								Path:        "Artist/Album/Disc1/01.mp3",
								HasCoverArt: false,
								DiscNumber:  1,
							},
							{
								Path:        "Artist/Album/Disc2/01.mp3",
								HasCoverArt: true,
								DiscNumber:  2,
							},
						}
					})
					It("selects the file with cover art even if from a higher disc number", func() {
						album := mfs.ToAlbum()
						Expect(album.EmbedArtPath).To(Equal("Artist/Album/Disc2/01.mp3"))
					})
				})

				When("we have media files with path names that don't correlate with disc numbers", func() {
					BeforeEach(func() {
						mfs = MediaFiles{
							{
								Path:        "Artist/Album/file-z.mp3", // Path would be sorted last alphabetically
								HasCoverArt: true,
								DiscNumber:  1, // But it has lowest disc number
							},
							{
								Path:        "Artist/Album/file-a.mp3", // Path would be sorted first alphabetically
								HasCoverArt: true,
								DiscNumber:  2, // But it has higher disc number
							},
							{
								Path:        "Artist/Album/file-m.mp3",
								HasCoverArt: true,
								DiscNumber:  3,
							},
						}
					})
					It("selects the cover art from the lowest disc number regardless of path", func() {
						album := mfs.ToAlbum()
						Expect(album.EmbedArtPath).To(Equal("Artist/Album/file-z.mp3"))
					})
				})
			})
		})
	})

	Describe("ToM3U8", func() {
		It("returns header only for empty MediaFiles", func() {
			mfs = MediaFiles{}
			result := mfs.ToM3U8("My Playlist", false)
			Expect(result).To(Equal("#EXTM3U\n#PLAYLIST:My Playlist\n"))
		})

		DescribeTable("duration formatting",
			func(duration float32, expected string) {
				mfs = MediaFiles{{Title: "Song", Artist: "Artist", Duration: duration, Path: "song.mp3"}}
				result := mfs.ToM3U8("Test", false)
				Expect(result).To(ContainSubstring(expected))
			},
			Entry("zero duration", float32(0.0), "#EXTINF:0,"),
			Entry("whole number", float32(120.0), "#EXTINF:120,"),
			Entry("rounds 0.5 down", float32(180.5), "#EXTINF:180,"),
			Entry("rounds 0.6 up", float32(240.6), "#EXTINF:241,"),
		)

		Context("multiple tracks", func() {
			BeforeEach(func() {
				mfs = MediaFiles{
					{Title: "Song One", Artist: "Artist A", Duration: 120, Path: "a/song1.mp3", LibraryPath: "/music"},
					{Title: "Song Two", Artist: "Artist B", Duration: 241, Path: "b/song2.mp3", LibraryPath: "/music"},
					{Title: "Song with \"quotes\" & ampersands", Artist: "Artist with Ümläuts", Duration: 90, Path: "special/file.mp3", LibraryPath: "/música"},
				}
			})

			DescribeTable("generates correct output",
				func(absolutePaths bool, expectedContent string) {
					result := mfs.ToM3U8("Multi Track", absolutePaths)
					Expect(result).To(Equal(expectedContent))
				},
				Entry("relative paths",
					false,
					"#EXTM3U\n#PLAYLIST:Multi Track\n#EXTINF:120,Artist A - Song One\na/song1.mp3\n#EXTINF:241,Artist B - Song Two\nb/song2.mp3\n#EXTINF:90,Artist with Ümläuts - Song with \"quotes\" & ampersands\nspecial/file.mp3\n",
				),
				Entry("absolute paths",
					true,
					"#EXTM3U\n#PLAYLIST:Multi Track\n#EXTINF:120,Artist A - Song One\n/music/a/song1.mp3\n#EXTINF:241,Artist B - Song Two\n/music/b/song2.mp3\n#EXTINF:90,Artist with Ümläuts - Song with \"quotes\" & ampersands\n/música/special/file.mp3\n",
				),
				Entry("special characters",
					false,
					"#EXTM3U\n#PLAYLIST:Multi Track\n#EXTINF:120,Artist A - Song One\na/song1.mp3\n#EXTINF:241,Artist B - Song Two\nb/song2.mp3\n#EXTINF:90,Artist with Ümläuts - Song with \"quotes\" & ampersands\nspecial/file.mp3\n",
				),
			)
		})

		Context("path variations", func() {
			It("handles different path structures", func() {
				mfs = MediaFiles{
					{Title: "Root", Artist: "Artist", Duration: 60, Path: "song.mp3", LibraryPath: "/lib"},
					{Title: "Nested", Artist: "Artist", Duration: 60, Path: "deep/nested/song.mp3", LibraryPath: "/lib"},
				}

				relativeResult := mfs.ToM3U8("Test", false)
				Expect(relativeResult).To(ContainSubstring("song.mp3\n"))
				Expect(relativeResult).To(ContainSubstring("deep/nested/song.mp3\n"))

				absoluteResult := mfs.ToM3U8("Test", true)
				Expect(absoluteResult).To(ContainSubstring("/lib/song.mp3\n"))
				Expect(absoluteResult).To(ContainSubstring("/lib/deep/nested/song.mp3\n"))
			})
		})
	})
})

var _ = Describe("MediaFile", func() {
	BeforeEach(func() {
		DeferCleanup(configtest.SetupConfig())
		conf.Server.EnableMediaFileCoverArt = true
	})
	Describe(".CoverArtId()", func() {
		It("returns its own id if it HasCoverArt", func() {
			mf := MediaFile{ID: "111", AlbumID: "1", HasCoverArt: true}
			id := mf.CoverArtID()
			Expect(id.Kind).To(Equal(KindMediaFileArtwork))
			Expect(id.ID).To(Equal(mf.ID))
		})
		It("returns its album id if HasCoverArt is false", func() {
			mf := MediaFile{ID: "111", AlbumID: "1", HasCoverArt: false}
			id := mf.CoverArtID()
			Expect(id.Kind).To(Equal(KindAlbumArtwork))
			Expect(id.ID).To(Equal(mf.AlbumID))
		})
		It("returns its album id if EnableMediaFileCoverArt is disabled", func() {
			conf.Server.EnableMediaFileCoverArt = false
			mf := MediaFile{ID: "111", AlbumID: "1", HasCoverArt: true}
			id := mf.CoverArtID()
			Expect(id.Kind).To(Equal(KindAlbumArtwork))
			Expect(id.ID).To(Equal(mf.AlbumID))
		})
	})
})

func t(v string) time.Time {
	var timeFormats = []string{"2006-01-02", "2006-01-02 15:04", "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02T15:04", "2006-01-02 15:04:05.999999999 -0700 MST"}
	for _, f := range timeFormats {
		t, err := time.ParseInLocation(f, v, time.UTC)
		if err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
